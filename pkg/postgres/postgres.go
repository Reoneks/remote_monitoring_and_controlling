package postgres

import (
	"context"
	"errors"
	"fmt"
	"project/config"
	"project/structs"

	"github.com/golang-migrate/migrate/v4"
	mpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

type Postgres struct {
	db *gorm.DB

	dsn        string
	migrations string
}

func (p *Postgres) GetUserByPhone(ctx context.Context, phone string) (structs.User, error) {
	var result structs.User
	return result, p.db.Model(&result).Where("phone = ?", phone).First(&result).Error
}

func (p *Postgres) GetUserByID(ctx context.Context, userID string) (structs.User, error) {
	var result structs.User
	return result, p.db.Model(&result).Where("id = ?", userID).First(&result).Error
}

func (p *Postgres) CreateUser(ctx context.Context, user *structs.User) error {
	return p.db.Model(user).Create(user).Error
}

func (p *Postgres) ChangePassword(ctx context.Context, userID, password string) error {
	return p.db.Model(&structs.User{}).Where("id = ?", userID).Update("password", password).Error
}

func (p *Postgres) EnableOTP(ctx context.Context, userID, secret string) error {
	return p.db.Model(&structs.User{}).Where("id = ?", userID).Updates(map[string]any{
		"otp_enabled": true,
		"otp_secret":  secret,
	}).Error
}

func (p *Postgres) BindTelegramUser(ctx context.Context, userPhone string, telegramUserID int64) error {
	return p.db.Model(&structs.User{}).Where("phone = ?", userPhone).Update("telegram_user_id", telegramUserID).Error
}

func (p *Postgres) Start(ctx context.Context) (err error) {
	p.db, err = newDB(p.dsn, p.migrations)
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

func newDB(dsn, migrationsURL string) (*gorm.DB, error) {
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: glog.Default.LogMode(glog.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("Postgres.newDB gorm connection open error: %w", err)
	}

	if err := migrations(client, dsn, migrationsURL); err != nil {
		return nil, fmt.Errorf("Postgres.newDB gorm migrations error: %w", err)
	}

	return client, nil
}

func migrations(client *gorm.DB, dsn, migrationsURL string) error {
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

func NewPostgres(cfg *config.Config) *Postgres {
	postgres := new(Postgres)
	postgres.dsn = cfg.DSN
	postgres.migrations = cfg.MigrationURL
	return postgres
}
