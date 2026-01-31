package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nrnasyrova/fx-rate-service/internal/app"
)

type GetRateByReqIdHandler struct {
	service *app.RateService
}

type getRateByIdResponse struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Status    string    `json:"status"`
	QuoteE6   *int64    `json:"quote_e6,omitempty"`
	ErrorMsg  *string   `json:"error_message,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetRateByReqIdHandler(service *app.RateService) *GetRateByReqIdHandler {
	return &GetRateByReqIdHandler{service: service}
}

func (h *GetRateByReqIdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqID := r.URL.Query().Get("id")
	if reqID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	refreshReq, err := h.service.GetByReqId(r.Context(), reqID)
	if err != nil {
		http.Error(w, fmt.Sprintf("internal err: %v", err), http.StatusInternalServerError)
		return
	}

	var quoteE6 *int64
	if refreshReq.ValueE6 != nil {
		v := int64(*refreshReq.ValueE6)
		quoteE6 = &v
	}

	response := getRateByIdResponse{
		ID:        refreshReq.ID,
		From:      refreshReq.Pair.BaseCurrency().String(),
		To:        refreshReq.Pair.QuoteCurrency().String(),
		Status:    string(refreshReq.Status),
		QuoteE6:   quoteE6,
		ErrorMsg:  refreshReq.ErrorMsg,
		UpdatedAt: refreshReq.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
