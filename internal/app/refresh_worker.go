package app

import (
	"context"
	"log"

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
				errMsg := err.Error()
				if updateErr := rs.rateRefreshRepo.Update(ctx, task.id, nil, models.Error, &errMsg); updateErr != nil {
					log.Printf("refresh worker: update failure status for %s: %v", task.id, updateErr)
				}
				continue
			}

			quote := models.NewValueE6FromFloat(price)

			// TODO: use a single DB transaction for Upsert + Update.
			if err := rs.rateRepo.Upsert(ctx, task.pair, quote); err != nil {
				errMsg := err.Error()
				if updateErr := rs.rateRefreshRepo.Update(ctx, task.id, nil, models.Error, &errMsg); updateErr != nil {
					log.Printf("refresh worker: update failure status for %s: %v", task.id, updateErr)
				}
				continue
			}

			if err := rs.rateRefreshRepo.Update(ctx, task.id, &quote, models.Success, nil); err != nil {
				log.Printf("refresh worker: update success status for %s: %v", task.id, err)
			}

		case <-ctx.Done():
			return // Graceful shutdown
		}
	}
}
