package app

import (
	"context"
	"fmt"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RefreshProcessor interface {
	ProcessRefresh(ctx context.Context, id string, pair models.CurrencyPair) error
}

func (rs *RateService) ProcessRefresh(ctx context.Context, id string, pair models.CurrencyPair) (err error) {
	defer func() {
		if err != nil {
			rs.handleRefreshError(ctx, id, err)
		}
	}()

	price, err := rs.rateProvider.FetchRate(ctx, pair)
	if err != nil {
		return fmt.Errorf("provider fetch: %w", err)
	}

	valueE6 := models.NewValueE6FromFloat(price)

	err = rs.txManager.WithTx(ctx, func(ctx context.Context, rateRepo RateRepository, refreshRepo RefreshRequestRepository) error {
		if err := rateRepo.Upsert(ctx, pair, valueE6); err != nil {
			return err
		}
		return refreshRepo.Update(ctx, id, &valueE6, models.Success, nil)
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func (rs *RateService) handleRefreshError(ctx context.Context, id string, originalErr error) {
	errMsg := originalErr.Error()
	_ = rs.rateRefreshRepo.Update(ctx, id, nil, models.Error, &errMsg)
}
