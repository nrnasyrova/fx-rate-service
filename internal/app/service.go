package app

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/domain/models"
	"github.com/nrnasyrova/fx-rate-service/internal/domain/services"
)

type RateService struct {
	repo services.RateRepository
}

func NewRateService(repo services.RateRepository) *RateService {
	return &RateService{repo: repo}
}

func (h *RateService) GetLatest(ctx context.Context, pair models.CurrencyPair) (models.Rate, error) {
	return h.repo.GetLatest(ctx, pair)
}
