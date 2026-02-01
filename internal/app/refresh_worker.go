package app

import (
	"context"
	"log"
)

type RefreshWorker struct {
	processor RefreshProcessor
	tasks     <-chan refreshTask
}

func NewRefreshWorker(processor RefreshProcessor, tasks <-chan refreshTask) *RefreshWorker {
	return &RefreshWorker{processor: processor, tasks: tasks}
}

func (w *RefreshWorker) Run(ctx context.Context) {
	for {
		select {
		case task := <-w.tasks:
			if err := w.processor.ProcessRefresh(ctx, task.id, task.pair); err != nil {
				log.Printf("refresh worker: process task %s failed: %v", task.id, err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (rs *RateService) StartWorker(ctx context.Context) {
	NewRefreshWorker(rs, rs.refreshChan).Run(ctx)
}
