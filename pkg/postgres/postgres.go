package postgres

import (
	"context"
	"errors"
	"fmt"

	"remote_monitoring_and_controlling/config"

	"github.com/golang-migrate/migrate/v4"
	mpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Needs for correct migrations
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	glog "gorm.io/gorm/logger"
)

type Postgres struct {
	db  *gorm.DB
	cfg *config.PostgresConfig
}

func (p *Postgres) GetUserByPhone(ctx context.Context, phone string) (User, error) {
	var result User
	return result, p.db.WithContext(ctx).Model(&User{}).
		Select([]string{
			"users.id",
			"users.department",
			"users.position",
			"users.full_name",
			"users.foreign_id",
			"users.password",
			"users.otp_secret",
		}).
		Joins("INNER JOIN contact_info ON users.id = contact_info.user_id").
		Where("contact_info.phone = ?", phone).
		First(&result).Error
}

func (p *Postgres) CreateUser(ctx context.Context, user *User) error {
	return p.db.Model(user).WithContext(ctx).Create(user).Error
}

func (p *Postgres) AddContactInfo(ctx context.Context, contactInfo *ContactInfo) error {
	return p.db.Model(contactInfo).WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "phone"}},
		DoUpdates: clause.AssignmentColumns([]string{"type"}),
	}).Create(contactInfo).Error
}

func (p *Postgres) EnableOTP(ctx context.Context, userID, secret string) error {
	return p.db.Model(&User{}).WithContext(ctx).Where("id = ?", userID).Update("otp_secret", secret).Error
}

func (p *Postgres) DisableOTP(ctx context.Context, userID string) error {
	return p.db.Model(&User{}).WithContext(ctx).Where("id = ?", userID).Update("otp_secret", "").Error
}

func (p *Postgres) DeleteUser(ctx context.Context, userID string) error {
	return p.db.Model(&User{}).WithContext(ctx).Where("id = ?", userID).Delete(&User{}).Error
}

func (p *Postgres) Start(ctx context.Context) (err error) {
	p.db, err = newDB(p.cfg.DSN, p.cfg.MigrationURL)
	return
}

func (p *Postgres) Stop(ctx context.Context) error {
	sql, err := p.db.WithContext(ctx).DB()
	if err != nil {
		return fmt.Errorf("Postgres.Stop get database error: %w", err)
	}

	err = sql.Close()
	if err != nil {
		return fmt.Errorf("Postgres.Stop connection close error: %w", err)
	}

	return nil
}

func NewPostgres(cfg *config.Config) *Postgres {
	return &Postgres{
		cfg: &cfg.PostgresConfig,
	}
}

func newDB(dsn, migrationsURL string) (*gorm.DB, error) {
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: glog.Default.LogMode(glog.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("Postgres.newDB gorm connection open error: %w", err)
	}

	if err := migrations(client, migrationsURL); err != nil {
		return nil, fmt.Errorf("Postgres.newDB gorm migrations error: %w", err)
	}

	return client, nil
}

func migrations(client *gorm.DB, migrationsURL string) error {
	db, err := client.DB()
	if err != nil {
		return fmt.Errorf("Postgres.migrations failed to get sql connection: %w", err)
	}

	driver, err := mpostgres.WithInstance(db, &mpostgres.Config{})
	if err != nil {
		return fmt.Errorf("Postgres.migrations create migration instance error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("Postgres.migrations migration init error: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("Postgres.migrations migration error: %w", err)
	}

	return nil
}
