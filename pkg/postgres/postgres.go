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
	glog "gorm.io/gorm/logger"
)

type Postgres struct {
	db  *gorm.DB
	cfg *config.PostgresConfig
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
