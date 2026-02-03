package api

import (
	"net/http"
)

func NewRouter(svc RateService) http.Handler {
	mux := http.NewServeMux()

	latest := NewGetLatestRateHandler(svc)
	refresh := NewRefreshRateHandler(svc)
	getRefreshReqById := NewGetRefreshRateByReqIdHandler(svc)

	mux.Handle("GET /rates/latest", latest)
	mux.Handle("POST /refresh-requests", refresh)
	mux.Handle("GET /refresh-requests", getRefreshReqById)

	return mux
}
