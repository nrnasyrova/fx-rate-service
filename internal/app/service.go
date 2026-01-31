package app

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RateService struct {
	rateRepo        RateRepository
	rateRefreshRepo RateRefreshRepository
	rateProvider    RateProvider
	refreshChan     chan refreshTask
}

type refreshTask struct {
	id   string
	pair models.CurrencyPair
}

func NewRateService(rateRepo RateRepository, rateRefreshRepo RateRefreshRepository, rateProvider RateProvider) *RateService {
	return &RateService{
		rateRepo:        rateRepo,
		rateRefreshRepo: rateRefreshRepo,
		rateProvider:    rateProvider,
		refreshChan:     make(chan refreshTask, 100),
	}
}

func (rs *RateService) GetLatest(ctx context.Context, pair models.CurrencyPair) (models.Rate, error) {
	return rs.rateRepo.GetLatest(ctx, pair)
}

func (rs *RateService) RefreshRate(ctx context.Context, pair models.CurrencyPair) (string, error) {
	refreshID, created, err := rs.rateRefreshRepo.GetOrCreateRequest(ctx, pair)

	if err != nil {
		return "", err
	}

	if created {
		rs.refreshChan <- refreshTask{refreshID, pair}
	}

	return refreshID, nil
}

func (rs *RateService) GetByReqId(ctx context.Context, reqId string) (models.Rate, error) {
	return rs.rateRefreshRepo.Get(ctx, reqId)
}
