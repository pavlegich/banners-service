// Package config contains server configuration
// objects and methods.
package config

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config contains values of server flags and environments.
type Config struct {
	Address           string        `env:"ADDRESS" json:"address"`
	DSN               string        `env:"DATABASE_DSN" json:"database_dsn"`
	CleanupInterval   time.Duration `env:"CLEANUP_INTERVAL" json:"cleanup_interval"`
	DefaultExpiration time.Duration `env:"DEFAULT_EXPIRATION" json:"default_expiration"`
}

// NewConfig returns new server config.
func NewConfig(ctx context.Context) *Config {
	return &Config{}
}

// ParseFlags handles and processes flags and environments values
// when launching the server.
func (cfg *Config) ParseFlags(ctx context.Context) error {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "HTTP-server endpoint address host:port")
	flag.StringVar(&cfg.DSN, "d", "postgresql://localhost:5432/postgres", "URI (DSN) to database")
	flag.DurationVar(&cfg.CleanupInterval, "clean", time.Duration(10)*time.Minute, "Interval for cleaning the expired banners")
	flag.DurationVar(&cfg.DefaultExpiration, "exp", time.Duration(5)*time.Minute, "Default time of expiration for banners")

	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("ParseFlags: wrong environment values %w", err)
	}

	return nil
}
