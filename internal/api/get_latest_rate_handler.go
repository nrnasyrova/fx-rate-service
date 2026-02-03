package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type GetLatestRateHandler struct {
	service RateService
}

type getLatestRateResponse struct {
	Pair      string    `json:"pair"`
	ValueE6   int64     `json:"value_e6"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetLatestRateHandler(service RateService) *GetLatestRateHandler {
	return &GetLatestRateHandler{service: service}
}

func (h *GetLatestRateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pair := r.URL.Query().Get("pair")

	parsedPair, err := models.ParseCurrencyPair(pair)
	if err != nil {
		handleError(w, err)
		return
	}

	latest, err := h.service.GetLatest(r.Context(), parsedPair)
	if err != nil {
		handleError(w, err)
		return
	}

	response := getLatestRateResponse{
		Pair:      latest.Pair.String(),
		ValueE6:   int64(latest.Quote),
		UpdatedAt: latest.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
