package api

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RateService interface {
	GetLatest(ctx context.Context, pair models.CurrencyPair) (models.Rate, error)
	RefreshRate(ctx context.Context, pair models.CurrencyPair) (string, error)
	GetRefreshRequest(ctx context.Context, id string) (models.RefreshRequest, error)
}
