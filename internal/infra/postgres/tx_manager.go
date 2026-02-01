package postgres

import (
	"context"
	"database/sql"

	"github.com/nrnasyrova/fx-rate-service/internal/app"
)

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context, rateRepo app.RateRepository, refreshRepo app.RefreshRequestRepository) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	rateRepo := NewRateRepositoryTx(tx)
	refreshRepo := NewRefreshRateRepositoryTx(tx)

	if err := fn(ctx, rateRepo, refreshRepo); err != nil {
		return err
	}

	return tx.Commit()
}
