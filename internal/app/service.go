package app

import (
	"context"
	"log"

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

func NewRateService(rateRepo RateRepository, rateRefreshRepo RefreshRequestRepository, rateProvider RateProvider, txManager TxManager, queueSize int) *RateService {
	return &RateService{
		rateRepo:        rateRepo,
		rateRefreshRepo: rateRefreshRepo,
		rateProvider:    rateProvider,
		txManager:       txManager,
		refreshChan:     make(chan refreshTask, queueSize),
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
		select {
		case rs.refreshChan <- refreshTask{refreshReqID, pair}:
		default:
			errMsg := models.ErrServiceOverloaded.Error()
			err = rs.rateRefreshRepo.Update(ctx, refreshReqID, nil, models.Error, &errMsg)
			if err != nil {
				log.Printf("failed to update rate refresh req id %s, err: %s", refreshReqID, err)
			}
		}
	}

	return refreshReqID, nil
}

func (rs *RateService) GetRefreshRequest(ctx context.Context, id string) (models.RefreshRequest, error) {
	return rs.rateRefreshRepo.Get(ctx, id)
}
