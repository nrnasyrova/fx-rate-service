package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr       string
	RefreshSweeper struct {
		Interval   time.Duration
		StaleAfter time.Duration
	}
	RefreshWorker struct {
		NumWorkers int
		QueueSize  int
	}
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
	cfg.RefreshSweeper.Interval = getEnvDuration("REFRESH_SWEEP_INTERVAL", 30*time.Second)
	cfg.RefreshSweeper.StaleAfter = getEnvDuration("REFRESH_STALE_AFTER", 5*time.Minute)
	cfg.RefreshWorker.NumWorkers = getEnvInt("REFRESH_WORKER_NUM", 4)
	cfg.RefreshWorker.QueueSize = getEnvInt("REFRESH_WORKER_QUEUE_SIZE", 50)
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

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

func getEnvInt(key string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}
