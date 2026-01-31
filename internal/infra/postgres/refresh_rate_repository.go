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

func (r *RefreshRateRepository) Update(ctx context.Context, id string, quote models.ValueE6, status models.Status, errorMessage *string) error {
	const q = `
		UPDATE refresh_rate_requests
		SET value_e6 = $1, status = $2, error_message = $3, updated_at = NOW()
		WHERE id = $4;`

	_, err := r.db.ExecContext(ctx, q, quote, status, errorMessage, id)
	return err
}

func (r *RefreshRateRepository) Get(ctx context.Context, id string) (models.Rate, error) {
	const q = `
		SELECT base_currency, quote_currency, value_e6, updated_at 
		FROM refresh_rate_requests 
		WHERE id = $1`

	var (
		baseCurrency  string
		quoteCurrency string
		valueE6       int64
		updatedAt     time.Time
	)

	err := r.db.QueryRowContext(ctx, q, id).Scan(&baseCurrency, &quoteCurrency, &valueE6, &updatedAt)
	if err != nil {
		return models.Rate{}, err

	}

	pair, err := models.NewCurrencyPair(baseCurrency, quoteCurrency)
	if err != nil {
		return models.Rate{}, err
	}

	return models.Rate{
		Pair:      pair,
		Quote:     models.ValueE6(valueE6),
		UpdatedAt: updatedAt,
	}, nil
}
