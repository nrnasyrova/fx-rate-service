package main

import (
	"log"
	"net/http"

	"github.com/nrnasyrova/fx-rate-service/internal/api/get_latest_rate"
	"github.com/nrnasyrova/fx-rate-service/internal/app"
	"github.com/nrnasyrova/fx-rate-service/internal/infra"
)

func main() {
	//TODO configure DI if have time
	repo := infra.NewRateRepository()
	service := app.NewRateService(repo)
	handler := get_latest_rate.NewHandler(service)

	mux := http.NewServeMux()
	mux.Handle("/fx-rate/latest", handler)

	log.Printf("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
