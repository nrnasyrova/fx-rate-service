package app

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/domain/rate"
)

type RateRepository interface {
	GetLatest(ctx context.Context, pair rate.CurrencyPair) (rate.Rate, error)
}
