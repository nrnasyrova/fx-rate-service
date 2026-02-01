package app

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RateRepository interface {
	GetLatest(ctx context.Context, pair models.CurrencyPair) (models.Rate, error)
	Upsert(ctx context.Context, pair models.CurrencyPair, valueE6 models.ValueE6) error
}

type RefreshRequestRepository interface {
	Update(ctx context.Context, id string, valueE6 *models.ValueE6, status models.Status, errorMessage *string) error
	GetOrCreate(ctx context.Context, pair models.CurrencyPair) (string, bool, error)
	Get(ctx context.Context, id string) (models.RefreshRequest, error)
}

type RateProvider interface {
	FetchRate(ctx context.Context, pair models.CurrencyPair) (float64, error)
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, rateRepo RateRepository, refreshRepo RefreshRequestRepository) error) error
}
