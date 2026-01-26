package get_latest_rate

import (
	"encoding/json"
	"net/http"

	"github.com/nrnasyrova/fx-rate-service/internal/api"
	"github.com/nrnasyrova/fx-rate-service/internal/app"
	"github.com/nrnasyrova/fx-rate-service/internal/domain/models"
)

type Handler struct {
	service *app.RateService
}

func NewHandler(service *app.RateService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	pair, err := models.NewCurrencyPair(from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rate, err := h.service.GetLatest(r.Context(), pair)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(api.ToRateDTO(rate))
}
