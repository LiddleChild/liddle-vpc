package frankfurterapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

type Client struct {
	endpoint *url.URL
	client   *http.Client
}

func NewClient(endpoint *url.URL) Client {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	return Client{
		endpoint: endpoint,
		client:   client,
	}
}

func (c Client) RatesToTHB(ctx context.Context, currencies []string) (ExchangeRateResponse, error) {
	resp := ExchangeRateResponse{
		"THB": {
			Rate:      1,
			UpdatedAt: time.Now(),
		},
	}

	if len(currencies) == 0 {
		return resp, nil
	}

	quotes := make(map[string]struct{})
	for _, currency := range currencies {
		code := strings.ToUpper(strings.TrimSpace(currency))
		quotes[code] = struct{}{}
	}

	endpoint := c.endpoint.JoinPath("v2/rates")

	query := endpoint.Query()
	query.Set("base", "THB")
	query.Set("quotes", strings.Join(slices.Collect(maps.Keys(quotes)), ","))
	endpoint.RawQuery = query.Encode()

	respBody, err := c.request(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var rates []frankfurterRate
	if err := json.Unmarshal(respBody, &rates); err != nil {
		return nil, fmt.Errorf("decode Frankfurter response: %w", err)
	}

	requested := make(map[string]struct{}, len(quotes))
	for quote := range quotes {
		requested[quote] = struct{}{}
	}

	for _, rate := range rates {
		quote := strings.ToUpper(rate.Quote)
		if _, ok := quotes[quote]; !ok {
			continue
		}

		updatedAt, err := time.Parse(time.DateOnly, rate.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		resp[quote] = ExchangeRate{
			Rate:      rate.Rate,
			UpdatedAt: updatedAt,
		}
	}

	return resp, nil
}

func (c Client) request(ctx context.Context, endpoint *url.URL) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create Frankfurter request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send Frankfurter request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("Frankfurter request failed with status %d", resp.StatusCode)
	}

	buffers, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return buffers, nil
}
