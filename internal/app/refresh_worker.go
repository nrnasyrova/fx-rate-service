package app

import (
	"context"
	"log"
	"sync"
)

type RefreshWorker struct {
	processor RefreshProcessor
	tasks     <-chan refreshTask
}

func newRefreshWorker(processor RefreshProcessor, tasks <-chan refreshTask) *RefreshWorker {
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

func (rs *RateService) StartWorkerPool(ctx context.Context, numWorkers int, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rs.startWorker(ctx)
		}()
	}
	log.Printf("Started %d refresh workers", numWorkers)
}

func (rs *RateService) startWorker(ctx context.Context) {
	newRefreshWorker(rs, rs.refreshChan).Run(ctx)
}
