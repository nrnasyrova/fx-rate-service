package postgres

import (
	"context"
	"database/sql"

	"github.com/nrnasyrova/fx-rate-service/internal/domain/rate"
)

type RateRepository struct {
	db *sql.DB
}

func NewRateRepository(db *sql.DB) *RateRepository {
	return &RateRepository{
		db: db,
	}
}

func (r *RateRepository) GetLatest(ctx context.Context, pair rate.CurrencyPair) (rate.Rate, error) {
	//TODO after connecting to storage implement real logic
	return rate.Rate{
		Pair:        pair,
		Quote:       100,
		UpdatedAtMs: 1640995200000,
	}, nil
}
func (r *RateRepository) Upsert(ctx context.Context, pair rate.CurrencyPair, quote rate.QuoteE6) error {
	return nil
}
