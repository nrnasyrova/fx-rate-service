package api

import "github.com/nrnasyrova/fx-rate-service/internal/domain/models"

type RateDTO struct {
	From        string `json:"from"`
	To          string `json:"to"`
	QuoteE6     int64  `json:"quote_e6"`
	UpdatedAtMs int64  `json:"updated_at_ms"`
}

func ToRateDTO(r models.Rate) RateDTO {
	return RateDTO{
		From: r.Pair.From().String(),
		To:   r.Pair.To().String(),
		//TODO need to figure out how to return without loosing data, is int ok?
		QuoteE6: int64(r.Quote),
		//TODO should I return timestamp in UTC?
		UpdatedAtMs: r.UpdatedAtMs,
	}
}
