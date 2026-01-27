package main

import (
	"log"
	"net/http"

	"github.com/nrnasyrova/fx-rate-service/internal/api"
	"github.com/nrnasyrova/fx-rate-service/internal/app"
	"github.com/nrnasyrova/fx-rate-service/internal/infra/postgres"
)

func main() {
	//TODO configure DI if have time
	repo := postgres.NewRateRepository()
	service := app.NewRateService(repo)

	handler := api.NewRouter(service)

	log.Printf("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
