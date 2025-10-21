package lighter

type remoteOrderBookSubscriptionPayload struct {
	Type      string `json:"type"`
	Operation string `json:"-"`
	Coin      string `json:"-"`
	Channel   string `json:"channel"`
}

func (p remoteOrderBookSubscriptionPayload) GetChannel() string {
	return p.Type
}

func (p remoteOrderBookSubscriptionPayload) Key() string {
	// Deliberately exclude NSigFigs and Mantissa.
	return keyOrderBook(p.Coin)
}

type remoteTradesSubscriptionPayload struct {
	Type      string `json:"type"`
	Operation string `json:"-"`
	Coin      string `json:"-"`
	Channel   string `json:"channel"`
}

func (p remoteTradesSubscriptionPayload) GetChannel() string {
	return p.Type
}

func (p remoteTradesSubscriptionPayload) Key() string {
	// Deliberately exclude NSigFigs and Mantissa.
	return keyTrades(p.Coin)
}
