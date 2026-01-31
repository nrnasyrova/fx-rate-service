package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	HTTPAddr string
	Postgres struct {
		DSN string
	}
	Migrations struct {
		Dir string
	}
	RateProvider struct {
		BaseUrl     string
		AccessToken string
	}
}

func FromEnv() (Config, error) {
	var cfg Config

	cfg.HTTPAddr = getEnv("HTTP_ADDR", ":8080")
	cfg.Postgres.DSN = os.Getenv("POSTGRES_DSN")
	if strings.TrimSpace(cfg.Postgres.DSN) == "" {
		return Config{}, errors.New("POSTGRES_DSN is required")
	}
	cfg.Migrations.Dir = getEnv("MIGRATIONS_DIR", "./migrations")
	cfg.RateProvider.BaseUrl = getEnv("RATE_PROVIDER_BASE_URL", "https://api.exchangeratesapi.io")
	cfg.RateProvider.AccessToken = os.Getenv("RATE_PROVIDER_ACCESS_TOKEN")
	if strings.TrimSpace(cfg.RateProvider.AccessToken) == "" {
		return Config{}, errors.New("RATE_PROVIDER_ACCESS_TOKEN is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}
