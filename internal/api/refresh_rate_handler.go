package api

import (
	"encoding/json"
	"net/http"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RefreshRateHandler struct {
	service RateService
}

type refreshRequest struct {
	Pair string `json:"pair"`
}

type refreshResponse struct {
	ID string `json:"id"`
}

func NewRefreshRateHandler(service RateService) *RefreshRateHandler {
	return &RefreshRateHandler{service: service}
}

func (h *RefreshRateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currencyPair, err := models.ParseCurrencyPair(req.Pair)
	if err != nil {
		handleError(w, err)
		return
	}

	id, err := h.service.RefreshRate(r.Context(), currencyPair)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(refreshResponse{
		ID: id,
	})
}
