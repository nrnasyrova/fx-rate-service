package app

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RateRepository interface {
	GetLatest(ctx context.Context, pair models.CurrencyPair) (models.Rate, error)
	Upsert(ctx context.Context, pair models.CurrencyPair, valueE6 models.ValueE6) error
}

type RateRefreshRepository interface {
	Update(ctx context.Context, id string, valueE6 *models.ValueE6, status models.Status, errorMessage *string) error
	GetOrCreateRequest(ctx context.Context, pair models.CurrencyPair) (string, bool, error)
	Get(ctx context.Context, id string) (models.RefreshRateRequest, error)
}

type RateProvider interface {
	FetchRate(ctx context.Context, pair models.CurrencyPair) (float64, error)
}
