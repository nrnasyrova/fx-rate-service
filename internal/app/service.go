package app

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RateService struct {
	rateRepo        RateRepository
	rateRefreshRepo RefreshRequestRepository
	rateProvider    RateProvider
	txManager       TxManager
	refreshChan     chan refreshTask
}

type refreshTask struct {
	id   string
	pair models.CurrencyPair
}

func NewRateService(rateRepo RateRepository, rateRefreshRepo RefreshRequestRepository, rateProvider RateProvider, txManager TxManager) *RateService {
	return &RateService{
		rateRepo:        rateRepo,
		rateRefreshRepo: rateRefreshRepo,
		rateProvider:    rateProvider,
		txManager:       txManager,
		refreshChan:     make(chan refreshTask, 1000), //TODO if fills up what to do?
	}
}

func (rs *RateService) GetLatest(ctx context.Context, pair models.CurrencyPair) (models.Rate, error) {
	return rs.rateRepo.GetLatest(ctx, pair)
}

func (rs *RateService) RefreshRate(ctx context.Context, pair models.CurrencyPair) (string, error) {
	refreshReqID, created, err := rs.rateRefreshRepo.GetOrCreate(ctx, pair)

	if err != nil {
		return "", err
	}

	if created {
		rs.refreshChan <- refreshTask{refreshReqID, pair}
	}

	return refreshReqID, nil
}

func (rs *RateService) GetRefreshRequest(ctx context.Context, id string) (models.RefreshRequest, error) {
	return rs.rateRefreshRepo.Get(ctx, id)
}
