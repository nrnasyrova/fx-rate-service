package postgres

import (
	"context"
	"database/sql"
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

func (r *RefreshRateRepository) GetOrCreateRequest(ctx context.Context, pair models.CurrencyPair) (string, bool, error) {

	const q = `
        INSERT INTO refresh_rate_requests (base_currency, quote_currency, status)
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
		UPDATE refresh_rate_requests
		SET value_e6 = $1, status = $2, error_message = $3, updated_at = NOW()
		WHERE id = $4;`

	_, err := r.db.ExecContext(ctx, q, quote, status, errorMessage, id)
	return err
}

func (r *RefreshRateRepository) Get(ctx context.Context, id string) (models.RefreshRateRequest, error) {
	const q = `
		SELECT id, base_currency, quote_currency, value_e6, status, created_at, COALESCE(updated_at, created_at), error_message
		FROM refresh_rate_requests 
		WHERE id = $1`

	var (
		reqID         string
		baseCurrency  string
		quoteCurrency string
		valueE6       sql.NullInt64
		status        string
		createdAt     time.Time
		updatedAt     time.Time
		errorMessage  sql.NullString
	)

	err := r.db.QueryRowContext(ctx, q, id).Scan(&reqID, &baseCurrency, &quoteCurrency, &valueE6, &status, &createdAt, &updatedAt, &errorMessage)
	if err != nil {
		return models.RefreshRateRequest{}, err

	}

	pair, err := models.NewCurrencyPair(baseCurrency, quoteCurrency)
	if err != nil {
		return models.RefreshRateRequest{}, err
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

	return models.RefreshRateRequest{
		ID:        reqID,
		Pair:      pair,
		ValueE6:   quote,
		Status:    models.Status(status),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ErrorMsg:  errMsg,
	}, nil
}
