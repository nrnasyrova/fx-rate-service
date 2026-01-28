package postgres

import (
	"context"
	"database/sql"

	"github.com/nrnasyrova/fx-rate-service/internal/domain/rate"
)

type RateRefreshRepository struct {
	db *sql.DB
}

func NewRateRefreshRepository(db *sql.DB) *RateRefreshRepository {
	return &RateRefreshRepository{
		db: db,
	}
}

func (r *RateRefreshRepository) GetOrCreateRequest(ctx context.Context, pair rate.CurrencyPair) (string, bool, error) {
	//INSERT INTO rate_refresh_requests (pair, status, created_at)
	//VALUES ($1, 'pending', now())
	//ON CONFLICT (pair) WHERE status = 'pending'
	//DO UPDATE SET pair = rate_refresh_requests.pair
	//RETURNING id;

	return "", false, nil
}

func (r *RateRefreshRepository) Update(ctx context.Context, id string, quote rate.QuoteE6) error {
	return nil
}

func (r *RateRefreshRepository) Upsert(ctx context.Context, pair rate.CurrencyPair, quote rate.QuoteE6) error {
	return nil
}
