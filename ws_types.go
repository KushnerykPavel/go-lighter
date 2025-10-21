package lighter

import (
	"encoding/json"
	"strings"
)

const (
	ChannelOrderBook string = "order_book"
	ChannelConnected string = "connected"
	ChannelTrades    string = "trade"
)

type wsMessage struct {
	MarketID  string          `json:"-"`
	Channel   string          `json:"channel"`
	Type      string          `json:"type"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
	OrderBook json.RawMessage `json:"order_book"`
}

func (w wsMessage) GetChannel() string {
	if w.Channel == "" {
		return w.Type
	}
	splitted := strings.Split(w.Channel, ":")
	if len(splitted) > 1 {
		return splitted[0]
	}
	return w.Channel
}

func (w wsMessage) GetMarketID() string {
	splitted := strings.Split(w.Channel, ":")
	if len(splitted) > 1 {
		return splitted[1]
	}
	return ""
}

type (
	OrderBook struct {
		MarketID  string  `json:"market_id"`
		Asks      []Level `json:"asks"`
		Bids      []Level `json:"bids"`
		Offset    int64   `json:"offset"`
		Timestamp int64   `json:"timestamp"`
	}

	Level struct {
		Price    string `json:"price"`
		Quantity string `json:"size"`
	}

	Trades struct {
		MarketID string  `json:"market_id"`
		Channel  string  `json:"channel"`
		Nonce    int     `json:"nonce"`
		Trades   []Trade `json:"trades"`
		Type     string  `json:"type"`
	}

	Trade struct {
		TradeID                          int    `json:"trade_id"`
		TxHash                           string `json:"tx_hash"`
		Type                             string `json:"type"`
		MarketID                         int    `json:"market_id"`
		Size                             string `json:"size"`
		Price                            string `json:"price"`
		UsdAmount                        string `json:"usd_amount"`
		AskId                            int64  `json:"ask_id"`
		BidId                            int64  `json:"bid_id"`
		AskAccountID                     int    `json:"ask_account_id"`
		BidAccountID                     int64  `json:"bid_account_id"`
		IsMakerAsk                       bool   `json:"is_maker_ask"`
		BlockHeight                      int    `json:"block_height"`
		Timestamp                        int64  `json:"timestamp"`
		TakerPositionSizeBefore          string `json:"taker_position_size_before"`
		TakerEntryQuoteBefore            string `json:"taker_entry_quote_before"`
		TakerInitialMarginFractionBefore int    `json:"taker_initial_margin_fraction_before"`
		MakerFee                         int    `json:"maker_fee"`
		MakerPositionSizeBefore          string `json:"maker_position_size_before"`
		MakerEntryQuoteBefore            string `json:"maker_entry_quote_before"`
		MakerInitialMarginFractionBefore int    `json:"maker_initial_margin_fraction_before"`
	}
)
