package app

import (
	"context"
	"log"
	"time"
)

type StaleRefreshSweeper struct {
	repo         StaleRefreshMarker
	interval     time.Duration
	staleAfter   time.Duration
	sweepNowFunc func() time.Time
}

func NewStaleRefreshSweeper(repo StaleRefreshMarker, interval, staleAfter time.Duration) *StaleRefreshSweeper {
	return &StaleRefreshSweeper{
		repo:         repo,
		interval:     interval,
		staleAfter:   staleAfter,
		sweepNowFunc: time.Now,
	}
}

func (s *StaleRefreshSweeper) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			staleBefore := s.sweepNowFunc().Add(-s.staleAfter)
			err := s.repo.MarkStaleProcessingRequest(ctx, staleBefore)
			if err != nil {
				log.Printf("refresh sweeper: %v", err)
				continue
			}
		}
	}
}
