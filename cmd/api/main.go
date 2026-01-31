package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
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

	defer db.Close()

	migrate(db, &cfg.Migrations.Dir)

	//TODO configure DI if have time
	rateRepo := postgres.NewRateRepository(db)
	refreshRepo := postgres.NewRefreshRateRepository(db)
	rateProvider := exchange_rates_api.NewClient(cfg.RateProvider.BaseUrl, cfg.RateProvider.AccessToken, 10*time.Second)
	service := app.NewRateService(rateRepo, refreshRepo, rateProvider)
	//TODO pass cancellable context
	go service.StartWorker(context.Background())

	handler := api.NewRouter(service)

	log.Printf("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
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
