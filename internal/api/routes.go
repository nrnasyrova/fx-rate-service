package api

import (
	"net/http"
)

func NewRouter(svc RateService) http.Handler {
	mux := http.NewServeMux()

	latest := NewGetLatestRateHandler(svc)
	refresh := NewRefreshRateHandler(svc)
	getRefreshReqById := NewGetRefreshRateByReqIdHandler(svc)

	mux.Handle("/rates/latest", latest)
	mux.Handle("/rates/refresh", refresh)
	mux.Handle("/refresh-requests", getRefreshReqById)

	return mux
}
