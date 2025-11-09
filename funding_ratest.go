package lighter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type fundingResponse struct {
	Code         int           `json:"code"`
	FundingRates []FundingRate `json:"funding_rates"`
}

type FundingRate struct {
	Symbol   string  `json:"symbol"`
	Exchange string  `json:"exchange"`
	Rate     float64 `json:"rate"`
}

func (c *Client) FundingRates(ctx context.Context) ([]FundingRate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURL+"funding-rates", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var fr fundingResponse
	if err := json.NewDecoder(resp.Body).Decode(&fr); err != nil {
		return nil, err
	}
	return fr.FundingRates, nil
}
