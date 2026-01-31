package exchange_rates_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type Client struct {
	baseURL     string
	accessToken string
	client      *http.Client
}

func NewClient(baseUrl string, accessToken string, timeout time.Duration) *Client {
	return &Client{
		baseURL:     baseUrl,
		accessToken: accessToken,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

type rateResponse struct {
	Success *bool `json:"success"`
	Error   *struct {
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error"`
	Timestamp *int64             `json:"timestamp"`
	Rates     map[string]float64 `json:"rates"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
}

func (c *Client) FetchRate(ctx context.Context, pair models.CurrencyPair) (float64, error) {
	u, err := url.Parse(c.baseURL + "/latest")
	if err != nil {
		return 0, fmt.Errorf("parse base url: %w", err)
	}

	q := u.Query()
	q.Set("base", strings.ToUpper(pair.BaseCurrency().String()))
	q.Set("symbols", strings.ToUpper(pair.QuoteCurrency().String()))
	q.Set("access_key", c.accessToken)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("provider status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed rateResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return 0, fmt.Errorf("decode: %w", err)
	}
	if parsed.Success != nil && !*parsed.Success {
		if parsed.Error != nil {
			return 0, fmt.Errorf("provider error: %s (%s)", parsed.Error.Type, parsed.Error.Info)
		}
		return 0, fmt.Errorf("provider error")
	}

	rate, ok := parsed.Rates[strings.ToUpper(pair.QuoteCurrency().String())]
	if !ok {
		return 0, fmt.Errorf("missing rate for %s", strings.ToUpper(pair.QuoteCurrency().String()))
	}

	return rate, nil
}
