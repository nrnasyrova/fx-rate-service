package main

import (
	"log"
	"net/http"

	"github.com/nrnasyrova/fx-rate-service/internal/api"
	"github.com/nrnasyrova/fx-rate-service/internal/app"
	"github.com/nrnasyrova/fx-rate-service/internal/infra/postgres"
	"github.com/nrnasyrova/fx-rate-service/internal/infra/provider/exchange_rates_api"
)

func main() {
	//TODO configure DI if have time
	rateRepo := postgres.NewRateRepository()
	refreshRepo := postgres.NewRateRefreshRepository()
	rateProvider := exchange_rates_api.NewClient()
	service := app.NewRateService(rateRepo, refreshRepo, rateProvider)

	handler := api.NewRouter(service)

	log.Printf("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
