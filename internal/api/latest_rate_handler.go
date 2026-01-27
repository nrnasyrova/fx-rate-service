package api

import (
	"encoding/json"
	"net/http"

	"github.com/nrnasyrova/fx-rate-service/internal/app"
	"github.com/nrnasyrova/fx-rate-service/internal/domain/rate"
)

type LatestRateHandler struct {
	service *app.RateService
}

type latestRateResponse struct {
	From        string `json:"from"`
	To          string `json:"to"`
	QuoteE6     int64  `json:"quote_e6"`
	UpdatedAtMs int64  `json:"updated_at_ms"`
}

func NewLatestRateHandler(service *app.RateService) *LatestRateHandler {
	return &LatestRateHandler{service: service}
}

func (h *LatestRateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	pair, err := rate.NewCurrencyPair(from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	latest, err := h.service.GetLatest(r.Context(), pair)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	response := latestRateResponse{
		From:        latest.Pair.From().String(),
		To:          latest.Pair.To().String(),
		QuoteE6:     int64(latest.Quote),
		UpdatedAtMs: latest.UpdatedAtMs,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
