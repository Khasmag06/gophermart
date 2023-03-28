package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	//DatabaseURI          string `env:"DATABASE_URI"`
	DatabaseURI string `env:"DATABASE_URI" envDefault:"postgres://localhost:5432/postgres?sslmode=disable"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "Server address")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "Accrual system address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "Database URI")
	flag.Parse()
	return &cfg, err

}
