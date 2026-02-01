package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type LatestRateHandler struct {
	service RateService
}

type latestRateResponse struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	QuoteE6   int64     `json:"quote_e6"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewLatestRateHandler(service RateService) *LatestRateHandler {
	return &LatestRateHandler{service: service}
}

func (h *LatestRateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	pair, err := models.NewCurrencyPair(from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	latest, err := h.service.GetLatest(r.Context(), pair)
	if err != nil {
		http.Error(w, fmt.Sprintf("internal err: %v", err), http.StatusInternalServerError)
		return
	}

	response := latestRateResponse{
		From:      latest.Pair.BaseCurrency().String(),
		To:        latest.Pair.QuoteCurrency().String(),
		QuoteE6:   int64(latest.Quote),
		UpdatedAt: latest.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
