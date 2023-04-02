package config

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppConfig
	PostgresConfig
}

type AppConfig struct {
	AppAddr         string `envconfig:"APP_SERVER_ADDRESS"`
	FrontendAddr    string `envconfig:"FRONTEND_ADDRESS" required:"true"`
	JWTAccessSecret string `envconfig:"JWT_ACCESS_SECRET" required:"true"`
}

type PostgresConfig struct {
	DSN          string `envconfig:"DB_DSN" required:"true"`
	MigrationURL string `envconfig:"DB_MIGRATION_URL" default:"file://migrations"`
}

var (
	once   sync.Once
	config *Config
)

func Get() (*Config, error) {
	var err error
	once.Do(func() {
		var cfg Config
		_ = godotenv.Load(".env")

		if err = envconfig.Process("", &cfg); err != nil {
			return
		}

		config = &cfg
	})

	return config, err
}
