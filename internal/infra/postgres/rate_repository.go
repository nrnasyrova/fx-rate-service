package postgres

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/domain/rate"
)

type RateRepository struct {
}

func NewRateRepository() *RateRepository {
	return &RateRepository{}
}

func (r *RateRepository) GetLatest(ctx context.Context, pair rate.CurrencyPair) (rate.Rate, error) {
	//TODO after connecting to storage implement real logic
	return rate.Rate{
		Pair:        pair,
		Quote:       100,
		UpdatedAtMs: 1640995200000,
	}, nil
}
