package api

import (
	"errors"
	"net/http"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidCurrencyCode) || errors.Is(err, models.ErrSameCurrencies) || errors.Is(err, models.ErrInvalidPairFormat):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, models.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
