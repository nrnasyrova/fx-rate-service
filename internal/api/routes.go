package api

import (
	"net/http"

	"github.com/nrnasyrova/fx-rate-service/internal/app"
)

func NewRouter(svc *app.RateService) http.Handler {
	mux := http.NewServeMux()

	latest := NewLatestRateHandler(svc)
	refresh := NewRefreshRateHandler(svc)

	mux.Handle("/rates/latest", latest)
	mux.Handle("/rates/refresh", refresh)

	return mux
}
