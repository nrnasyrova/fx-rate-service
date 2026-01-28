package exchange_rates_api

import (
	"context"

	"github.com/nrnasyrova/fx-rate-service/internal/domain/rate"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) FetchRate(ctx context.Context, pair rate.CurrencyPair) (float64, error) {
	return 0, nil
}
