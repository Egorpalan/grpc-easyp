package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env            string `envconfig:"ENV" default:"local"`
	PostgresConfig PostgresConfig
	ServerConfig   ServerConfig
}

type PostgresConfig struct {
	PostgresHost     string `envconfig:"POSTGRES_HOST" default:"localhost" required:"true"`
	PostgresPort     string `envconfig:"POSTGRES_PORT" default:"5432" required:"true"`
	PostgresDB       string `envconfig:"POSTGRES_DB" default:"notes_db" required:"true"`
	PostgresUser     string `envconfig:"POSTGRES_USER" default:"notes_user" required:"true"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:"notes_password" required:"true"`
	PostgresMaxConn  int32  `envconfig:"POSTGRES_MAX_CONN" default:"10"`
	PostgresSslMode  string `envconfig:"POSTGRES_SSL_MODE" default:"disable"`
}

type ServerConfig struct {
	GRPCPort string `envconfig:"GRPC_PORT" default:"50051"`
	GRPCHost string `envconfig:"GRPC_HOST" default:"0.0.0.0"`
}

func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
		c.PostgresSslMode,
	)
}

func Load() (*Config, error) {
	cfg := &Config{}

	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
