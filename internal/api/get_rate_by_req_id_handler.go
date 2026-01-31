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
	Id        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	QuoteE6   int64     `json:"quote_e6"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetRateByReqIdHandler(service *app.RateService) *GetRateByReqIdHandler {
	return &GetRateByReqIdHandler{service: service}
}

func (h *GetRateByReqIdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqId := r.URL.Query().Get("id")

	rate, err := h.service.GetByReqId(r.Context(), reqId)
	if err != nil {
		http.Error(w, fmt.Sprintf("internal err: %v", err), http.StatusInternalServerError)
		return
	}

	response := getRateByIdResponse{
		Id:        reqId,
		From:      rate.Pair.BaseCurrency().String(),
		To:        rate.Pair.QuoteCurrency().String(),
		QuoteE6:   int64(rate.Quote),
		UpdatedAt: rate.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
