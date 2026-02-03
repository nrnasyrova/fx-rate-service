package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/nrnasyrova/fx-rate-service/internal/api"
	"github.com/nrnasyrova/fx-rate-service/internal/app"
	"github.com/nrnasyrova/fx-rate-service/internal/config"
	"github.com/nrnasyrova/fx-rate-service/internal/infra/exchange_rates_api"
	"github.com/nrnasyrova/fx-rate-service/internal/infra/postgres"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer closeDB(db)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup
	service := initRateService(cfg, db, &wg, ctx)

	router := api.NewRouter(service)
	if err := runHTTPServer(ctx, cfg.HTTPAddr, router); err != nil {
		log.Printf("server error: %v", err)
		stop()
	}

	wg.Wait()
}

func loadConfig() (*config.Config, error) {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Printf("env file load: %v", err)
		}
	}

	cfg, err := config.FromEnv()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return &cfg, nil
}

func initRateService(cfg *config.Config, db *sql.DB, wg *sync.WaitGroup, ctx context.Context) *app.RateService {
	rateRepo := postgres.NewRateRepository(db)
	refreshRepo := postgres.NewRefreshRateRepository(db)
	txManager := postgres.NewTxManager(db)
	rateProvider := exchange_rates_api.NewClient(
		cfg.RateProvider.BaseUrl,
		cfg.RateProvider.AccessToken,
		10*time.Second,
	)

	service := app.NewRateService(
		rateRepo,
		refreshRepo,
		rateProvider,
		txManager,
		cfg.RefreshWorker.QueueSize,
	)

	service.StartWorkerPool(ctx, cfg.RefreshWorker.NumWorkers, wg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		sweeper := app.NewStaleRefreshSweeper(refreshRepo, cfg.RefreshSweeper.Interval, cfg.RefreshSweeper.StaleAfter)
		sweeper.Run(ctx)
	}()

	return service
}

func initDB(cfg *config.Config) (db *sql.DB, err error) {
	db, err = postgres.Open(cfg.Postgres.DSN)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	defer func() {
		if err != nil && db != nil {
			_ = db.Close()
		}
	}()

	err = goose.SetDialect("postgres")
	if err != nil {
		return nil, fmt.Errorf("goose dialect: %w", err)
	}

	err = goose.Up(db, cfg.Migrations.Dir)
	if err != nil {
		return nil, fmt.Errorf("goose up: %w", err)
	}

	log.Printf("migrations applied")
	return db, nil
}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("db close: %v", err)
	}
}

func runHTTPServer(ctx context.Context, addr string, handler http.Handler) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Printf("Shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)

	case e := <-serverErr:
		if errors.Is(e, http.ErrServerClosed) {
			return nil
		}
		return e
	}
}
