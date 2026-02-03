package exchange_rates_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nrnasyrova/fx-rate-service/internal/models"
)

type Client struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
}

func NewClient(baseURL, accessToken string, timeout time.Duration) *Client {
	return &Client{
		baseURL:     strings.TrimSuffix(baseURL, "/"),
		accessToken: accessToken,
		httpClient:  &http.Client{Timeout: timeout},
	}
}

type rateResponse struct {
	Success bool               `json:"success"`
	Rates   map[string]float64 `json:"rates"`
	Error   struct {
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error"`
}

func (c *Client) FetchRate(ctx context.Context, pair models.CurrencyPair) (float64, error) {
	reqURL := c.buildLatestRateURL(pair)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	var resp rateResponse
	if err := c.execute(req, &resp); err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, fmt.Errorf("api error: %s - %s", resp.Error.Type, resp.Error.Info)
	}

	rate, ok := resp.Rates[pair.QuoteCurrency().String()]
	if !ok {
		return 0, fmt.Errorf("rate for %s not found", pair.QuoteCurrency().String())
	}

	return rate, nil
}

func (c *Client) buildLatestRateURL(pair models.CurrencyPair) string {
	params := url.Values{}
	params.Set("base", pair.BaseCurrency().String())
	params.Set("symbols", pair.QuoteCurrency().String())
	params.Set("access_key", c.accessToken)

	return fmt.Sprintf("%s%s?%s", c.baseURL, "/latest", params.Encode())
}

func (c *Client) execute(req *http.Request, target interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http execute: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
