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
}

func FromEnv() (Config, error) {
	var cfg Config

	cfg.HTTPAddr = getEnv("HTTP_ADDR", ":8080")
	cfg.Postgres.DSN = os.Getenv("POSTGRES_DSN")
	if strings.TrimSpace(cfg.Postgres.DSN) == "" {
		return Config{}, errors.New("POSTGRES_DSN is required")
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
