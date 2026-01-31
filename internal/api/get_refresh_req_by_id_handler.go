package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nrnasyrova/fx-rate-service/internal/app"
)

type GetRefreshReqByIdHandler struct {
	service *app.RateService
}

type getRefreshReqByIdResponse struct {
	ID        string    `json:"id"`
	Pair      string    `json:"pair"`
	Status    string    `json:"status"`
	QuoteE6   *int64    `json:"quote_e6,omitempty"`
	ErrorMsg  *string   `json:"error_message,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetRefreshRateByReqIdHandler(service *app.RateService) *GetRefreshReqByIdHandler {
	return &GetRefreshReqByIdHandler{service: service}
}

func (h *GetRefreshReqByIdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqID := r.URL.Query().Get("id")
	if reqID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	err := uuid.Validate(reqID)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	refreshReq, err := h.service.GetRefreshRequest(r.Context(), reqID)
	if err != nil {
		handleError(w, err)
		return
	}

	var quoteE6 *int64
	if refreshReq.ValueE6 != nil {
		v := int64(*refreshReq.ValueE6)
		quoteE6 = &v
	}

	response := getRefreshReqByIdResponse{
		ID:        refreshReq.ID,
		Pair:      refreshReq.Pair.String(),
		Status:    refreshReq.Status.String(),
		QuoteE6:   quoteE6,
		ErrorMsg:  refreshReq.ErrorMsg,
		UpdatedAt: refreshReq.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
