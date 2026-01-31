package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type RateRepository struct {
	db *sql.DB
}

func NewRateRepository(db *sql.DB) *RateRepository {
	return &RateRepository{
		db: db,
	}
}

func (r *RateRepository) GetLatest(ctx context.Context, pair models.CurrencyPair) (models.Rate, error) {
	const q = `
			SELECT base_currency, quote_currency, value_e6, updated_at 
			FROM rates 
			WHERE base_currency = $1 and quote_currency = $2
			LIMIT 1`

	var (
		baseCurrency  string
		quoteCurrency string
		valueE6       int64
		updatedAt     time.Time
	)

	err := r.db.QueryRowContext(ctx, q, pair.BaseCurrency(), pair.QuoteCurrency()).Scan(&baseCurrency, &quoteCurrency, &valueE6, &updatedAt)

	if err != nil {
		return models.Rate{}, err
	}

	pair, err = models.NewCurrencyPair(baseCurrency, quoteCurrency)

	if err != nil {
		return models.Rate{}, err
	}

	return models.Rate{
		Pair:      pair,
		Quote:     models.ValueE6(valueE6),
		UpdatedAt: updatedAt,
	}, nil
}

func (r *RateRepository) Upsert(ctx context.Context, pair models.CurrencyPair, quote models.ValueE6) error {
	const q = `
			INSERT INTO rates (base_currency, quote_currency, value_e6) 
			VALUES ($1, $2, $3)
			ON CONFLICT (base_currency, quote_currency)
			DO UPDATE SET value_e6 = $3, updated_at = NOW();`

	_, err := r.db.ExecContext(ctx, q, pair.BaseCurrency(), pair.QuoteCurrency(), quote)

	return err
}
