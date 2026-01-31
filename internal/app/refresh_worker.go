package app

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RefreshWorker struct {
}

func (rs *RateService) StartWorker(ctx context.Context) {
	for {
		select {
		case task := <-rs.refreshChan:
			price, err := rs.rateProvider.FetchRate(ctx, task.pair)
			if err != nil {

			}

			//TODO introduce transaction for upddting and inserting
			err = rs.rateRefreshRepo.Update(ctx, task.id, models.NewValueE6FromFloat(price), models.Success, nil)
			if err != nil {

			}

			err = rs.rateRepo.Upsert(ctx, task.pair, models.NewValueE6FromFloat(price))
			if err != nil {

			}

		case <-ctx.Done():
			return // Graceful shutdown
		}
	}
}
