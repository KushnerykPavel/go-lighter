package lighter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OrderBookDetail struct {
	Symbol                       string  `json:"symbol"`
	MarketId                     int     `json:"market_id"`
	Status                       string  `json:"status"`
	TakerFee                     string  `json:"taker_fee"`
	MakerFee                     string  `json:"maker_fee"`
	LiquidationFee               string  `json:"liquidation_fee"`
	MinBaseAmount                string  `json:"min_base_amount"`
	MinQuoteAmount               string  `json:"min_quote_amount"`
	SupportedSizeDecimals        int     `json:"supported_size_decimals"`
	SupportedPriceDecimals       int     `json:"supported_price_decimals"`
	SupportedQuoteDecimals       int     `json:"supported_quote_decimals"`
	SizeDecimals                 int     `json:"size_decimals"`
	PriceDecimals                int     `json:"price_decimals"`
	QuoteMultiplier              int     `json:"quote_multiplier"`
	DefaultInitialMarginFraction int     `json:"default_initial_margin_fraction"`
	MinInitialMarginFraction     int     `json:"min_initial_margin_fraction"`
	MaintenanceMarginFraction    int     `json:"maintenance_margin_fraction"`
	CloseoutMarginFraction       int     `json:"closeout_margin_fraction"`
	LastTradePrice               float64 `json:"last_trade_price"`
	DailyTradesCount             int     `json:"daily_trades_count"`
	DailyBaseTokenVolume         float64 `json:"daily_base_token_volume"`
	DailyQuoteTokenVolume        float64 `json:"daily_quote_token_volume"`
	DailyPriceLow                float64 `json:"daily_price_low"`
	DailyPriceHigh               float64 `json:"daily_price_high"`
	DailyPriceChange             float64 `json:"daily_price_change"`
	OpenInterest                 float64 `json:"open_interest"`
}

type orderBookDetailsResponse struct {
	Code             int               `json:"code"`
	OrderBookDetails []OrderBookDetail `json:"order_book_details"`
}

func (c *Client) OrderBookDetails(ctx context.Context) ([]OrderBookDetail, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURL+"orderBookDetails", nil)
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

	var fr orderBookDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&fr); err != nil {
		return nil, err
	}
	return fr.OrderBookDetails, nil
}
