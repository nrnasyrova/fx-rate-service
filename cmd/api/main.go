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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := postgres.Open(cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}

	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("db close: %v", err)
		}
	}(db)

	migrate(db, &cfg.Migrations.Dir)

	//TODO configure DI if have time
	rateRepo := postgres.NewRateRepository(db)
	refreshRepo := postgres.NewRefreshRateRepository(db)
	txManager := postgres.NewTxManager(db)
	rateProvider := exchange_rates_api.NewClient(cfg.RateProvider.BaseUrl, cfg.RateProvider.AccessToken, 10*time.Second)
	service := app.NewRateService(rateRepo, refreshRepo, rateProvider, txManager)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		service.StartWorker(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sweeper := app.NewStaleRefreshSweeper(refreshRepo, cfg.RefreshSweeper.Interval, cfg.RefreshSweeper.StaleAfter)
		sweeper.Run(ctx)
	}()

	router := api.NewRouter(service)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	srvErr := make(chan error, 1)
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Printf("Shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown: %v", err)
		}

		if err := <-srvErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
		}
	case err := <-srvErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
		}
		stop()
	}

	wg.Wait()
}

func migrate(db *sql.DB, dir *string) {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	if err := goose.Up(db, *dir); err != nil {
		log.Fatal(err)
	}

	fmt.Println("migrations applied")
}
