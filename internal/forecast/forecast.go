package forecast

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"alpineworks.io/rfc9457"
)

type ForecastClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewForecastClient(httpClient *http.Client, baseURL string, apiKey string) *ForecastClient {
	return &ForecastClient{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

type ForecastSummaryResponse struct {
	Summary     string `json:"summary"`
	Icon        string `json:"icon"`
	LastUpdated string `json:"last_updated"`
}

func (fc *ForecastClient) GetSummary(ctx context.Context) (*ForecastSummaryResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/forecast/summary", fc.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Api-Key", fc.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := fc.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body %w", err)
		}

		if body == nil {
			return nil, fmt.Errorf("%d empty body", resp.StatusCode)
		}

		problem, err := rfc9457.FromJSON(string(body))
		if err != nil {
			return nil, fmt.Errorf("failed to parse problem: %w body: %s code: %d", err, body, resp.StatusCode)
		}
		return nil, fmt.Errorf("problem: %v", problem)
	}

	var summary ForecastSummaryResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	err = json.Unmarshal(body, &summary)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return &summary, nil
}
