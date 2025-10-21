package lighter

type remoteProcessor interface {
	SetMarketID(marketID string)
	Key() string
}

func (c *OrderBook) SetMarketID(marketID string) {
	c.MarketID = marketID
}

func (c *Trades) SetMarketID(marketID string) {
	c.MarketID = marketID
}
