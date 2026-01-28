package app

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/domain/rate"
)

type RateRepository interface {
	GetLatest(ctx context.Context, pair rate.CurrencyPair) (rate.Rate, error)
	Upsert(ctx context.Context, pair rate.CurrencyPair, quote rate.QuoteE6) error
}

type RateRefreshRepository interface {
	Update(ctx context.Context, id string, quote rate.QuoteE6) error
	GetOrCreateRequest(ctx context.Context, pair rate.CurrencyPair) (string, bool, error)
}

type RateProvider interface {
	FetchRate(ctx context.Context, pair rate.CurrencyPair) (float64, error)
}
