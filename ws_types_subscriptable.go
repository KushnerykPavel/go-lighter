package lighter

type subscriptable interface {
	Key() string
}

func (c *OrderBook) Key() string {
	return keyOrderBook(c.MarketID)
}

func (c *Trades) Key() string {
	return keyTrades(c.MarketID)
}
