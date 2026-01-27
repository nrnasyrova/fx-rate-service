package app

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/domain/rate"
)

type RateService struct {
	repo RateRepository
}

func NewRateService(repo RateRepository) *RateService {
	return &RateService{repo: repo}
}

func (rs *RateService) GetLatest(ctx context.Context, pair rate.CurrencyPair) (rate.Rate, error) {
	return rs.repo.GetLatest(ctx, pair)
}

func (rs *RateService) RefreshRate(ctx context.Context, pair rate.CurrencyPair) (string, error) {
	return "", nil
}
