package app

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RefreshProcessor interface {
	ProcessRefresh(ctx context.Context, id string, pair models.CurrencyPair) error
}

func (rs *RateService) ProcessRefresh(ctx context.Context, id string, pair models.CurrencyPair) error {
	price, err := rs.rateProvider.FetchRate(ctx, pair)

	if err != nil {
		errMsg := err.Error()
		if updateErr := rs.rateRefreshRepo.Update(ctx, id, nil, models.Error, &errMsg); updateErr != nil {
			return updateErr
		}
		return nil
	}

	valueE6 := models.NewValueE6FromFloat(price)

	if err := rs.txManager.WithTx(ctx, func(ctx context.Context, rateRepo RateRepository, refreshRepo RefreshRequestRepository) error {
		if err := rateRepo.Upsert(ctx, pair, valueE6); err != nil {
			return err
		}
		return refreshRepo.Update(ctx, id, &valueE6, models.Success, nil)
	}); err != nil {
		errMsg := err.Error()
		if updateErr := rs.rateRefreshRepo.Update(ctx, id, nil, models.Error, &errMsg); updateErr != nil {
			return updateErr
		}
		return nil
	}

	return nil
}
