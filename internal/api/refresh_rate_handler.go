package api

import (
	"encoding/json"
	"net/http"

	"github.com/nrnasyrova/fx-rate-service/internal/app"
	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RefreshRateHandler struct {
	service *app.RateService
}

type refreshRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type refreshResponse struct {
	ID string `json:"id"`
}

func NewRefreshRateHandler(service *app.RateService) *RefreshRateHandler {
	return &RefreshRateHandler{service: service}
}

func (h *RefreshRateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currencyPair, err := models.NewCurrencyPair(req.From, req.To)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.RefreshRate(r.Context(), currencyPair)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(refreshResponse{
		ID: id,
	})
}
