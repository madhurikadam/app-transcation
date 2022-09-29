package postgres

import (
	"fmt"
	"time"
)

type Config struct {
	Database          string        `envconfig:"POSTGRES_DATABASE"`
	Host              string        `envconfig:"POSTGRES_HOST"`
	Password          string        `envconfig:"POSTGRES_PASSWORD"`
	Port              int           `envconfig:"POSTGRES_PORT"`
	User              string        `envconfig:"POSTGRES_USER"`
	SSLMode           string        `envconfig:"POSTGRES_SSL_MODE" default:"require"`
	MaxConnection     int32         `envconfig:"POSTGRES_MAX_CONNECTION" default:"25"`
	ConnectionTimeout int           `envconfig:"POSTRES_CONNECTION_TIMEOUT" default:"5"`
	AcquireTimeout    time.Duration `envconfig:"POSTGRES_ACQUIRE_TIMEOUT" default:"1s"`
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		// "postgres://username:password@localhost:5432/database_name"
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&connect_timeout=%d",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
		c.ConnectionTimeout,
	)
}
