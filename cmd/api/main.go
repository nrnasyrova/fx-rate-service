package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/nrnasyrova/fx-rate-service/internal/api"
	"github.com/nrnasyrova/fx-rate-service/internal/app"
	"github.com/nrnasyrova/fx-rate-service/internal/config"
	"github.com/nrnasyrova/fx-rate-service/internal/infra/exchange_rates_api"
	"github.com/nrnasyrova/fx-rate-service/internal/infra/postgres"
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

	//TODO configure DI if have time
	rateRepo := postgres.NewRateRepository(db)
	refreshRepo := postgres.NewRateRefreshRepository(db)
	rateProvider := exchange_rates_api.NewClient()
	service := app.NewRateService(rateRepo, refreshRepo, rateProvider)

	handler := api.NewRouter(service)

	log.Printf("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
