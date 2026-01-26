package infra

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/domain/models"
)

type RateRepository struct {
}

func NewRateRepository() *RateRepository {
	return &RateRepository{}
}

func (r *RateRepository) GetLatest(ctx context.Context, pair models.CurrencyPair) (models.Rate, error) {
	//TODO after connecting to storage implement real logic
	return models.Rate{
		Pair:        pair,
		Quote:       100,
		UpdatedAtMs: 1640995200000,
	}, nil
}
