package services

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/domain/models"
)

type RateRepository interface {
	GetLatest(ctx context.Context, pair models.CurrencyPair) (models.Rate, error)
}
