package finnhub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://finnhub.io/api/v1"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Earning represents a single earnings entry from Finnhub.
type Earning struct {
	Date            string  `json:"date"`
	Symbol          string  `json:"symbol"`
	Hour            string  `json:"hour"`             // bmo, amc, dmh
	Quarter         int     `json:"quarter"`           // fiscal quarter
	Year            int     `json:"year"`              // fiscal year
	EPSEstimate     *float64 `json:"epsEstimate"`
	EPSActual       *float64 `json:"epsActual"`
	RevenueEstimate *float64 `json:"revenueEstimate"`
	RevenueActual   *float64 `json:"revenueActual"`
}

type earningsResponse struct {
	EarningsCalendar []Earning `json:"earningsCalendar"`
}

// EarningsCalendar fetches earnings between from and to dates (YYYY-MM-DD).
// If symbol is non-empty, results are filtered to that single symbol server-side.
func (c *Client) EarningsCalendar(from, to, symbol string) ([]Earning, error) {
	reqURL := fmt.Sprintf("%s/calendar/earnings?from=%s&to=%s&token=%s", baseURL,
		url.QueryEscape(from), url.QueryEscape(to), url.QueryEscape(c.apiKey))
	if symbol != "" {
		reqURL += "&symbol=" + url.QueryEscape(symbol)
	}

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("finnhub request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("finnhub rate limit exceeded (free tier: 60 calls/min)")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("finnhub returned status %d: %s", resp.StatusCode, string(body))
	}

	var result earningsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode finnhub response: %w", err)
	}

	return result.EarningsCalendar, nil
}
