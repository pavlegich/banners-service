// Package config contains server configuration
// objects and methods.
package config

import (
	"context"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Config contains values of server flags and environments.
type Config struct {
	Address string `env:"ADDRESS" json:"address"`
	DSN     string `env:"DATABASE_DSN" json:"database_dsn"`
}

// NewConfig returns new server config.
func NewConfig(ctx context.Context) *Config {
	return &Config{}
}

// ParseFlags handles and processes flags and environments values
// when launching the server.
func (cfg *Config) ParseFlags(ctx context.Context) error {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "HTTP-server endpoint address host:port")
	flag.StringVar(&cfg.DSN, "d", "", "URI (DSN) to database")

	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("ParseFlags: wrong environment values %w", err)
	}

	return nil
}
