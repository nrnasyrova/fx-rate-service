package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RefreshRateRepository struct {
	db *sql.DB
}

func NewRefreshRateRepository(db *sql.DB) *RefreshRateRepository {
	return &RefreshRateRepository{
		db: db,
	}
}

func (r *RefreshRateRepository) GetOrCreate(ctx context.Context, pair models.CurrencyPair) (string, bool, error) {

	const q = `
        INSERT INTO refresh_requests (base_currency, quote_currency, status)
        VALUES ($1, $2, $3)
        ON CONFLICT (base_currency, quote_currency) WHERE status = 'processing'
        DO UPDATE SET updated_at = NOW()
        RETURNING id, (xmax = 0) AS created;
    `

	var id string
	var created bool
	err := r.db.QueryRowContext(ctx, q,
		pair.BaseCurrency().String(),
		pair.QuoteCurrency().String(),
		models.Processing,
	).Scan(&id, &created)

	if err != nil {
		return "", false, err
	}

	return id, created, nil
}

func (r *RefreshRateRepository) Update(ctx context.Context, id string, quote *models.ValueE6, status models.Status, errorMessage *string) error {
	const q = `
		UPDATE refresh_requests
		SET value_e6 = $1, status = $2, error_message = $3, updated_at = NOW()
		WHERE id = $4;`

	_, err := r.db.ExecContext(ctx, q, quote, status, errorMessage, id)
	return err
}

func (r *RefreshRateRepository) Get(ctx context.Context, id string) (models.RefreshRequest, error) {
	const q = `
		SELECT id, base_currency, quote_currency, value_e6, status, created_at, updated_at, error_message
		FROM refresh_requests 
		WHERE id = $1`

	var (
		ID            string
		baseCurrency  string
		quoteCurrency string
		valueE6       sql.NullInt64
		status        string
		createdAt     time.Time
		updatedAt     time.Time
		errorMessage  sql.NullString
	)

	err := r.db.QueryRowContext(ctx, q, id).Scan(&ID, &baseCurrency, &quoteCurrency, &valueE6, &status, &createdAt, &updatedAt, &errorMessage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.RefreshRequest{}, models.ErrNotFound
		}

		return models.RefreshRequest{}, fmt.Errorf("scanning refresh request: %w", err)
	}

	pair, err := models.NewCurrencyPair(baseCurrency, quoteCurrency)
	if err != nil {
		return models.RefreshRequest{}, err
	}

	var quote *models.ValueE6
	if valueE6.Valid {
		v := models.ValueE6(valueE6.Int64)
		quote = &v
	}

	var errMsg *string
	if errorMessage.Valid {
		v := errorMessage.String
		errMsg = &v
	}

	return models.RefreshRequest{
		ID:        ID,
		Pair:      pair,
		ValueE6:   quote,
		Status:    models.Status(status),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ErrorMsg:  errMsg,
	}, nil
}
