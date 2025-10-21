package lighter

import (
	"fmt"
)

type OrderBookSubscriptionParams struct {
	Coin string
}

func (w *WebsocketClient) OrderBook(
	params OrderBookSubscriptionParams,
	callback func(*OrderBook, error),
) (*Subscription, error) {
	remotePayload := remoteOrderBookSubscriptionPayload{
		Type:      "subscribe",
		Operation: ChannelOrderBook,
		Coin:      params.Coin,
		Channel:   fmt.Sprintf("%s/%s", ChannelOrderBook, params.Coin),
	}

	return w.subscribe(remotePayload, func(msg any) {
		orderbook, ok := msg.(*OrderBook)
		if !ok {
			callback(&OrderBook{}, fmt.Errorf("invalid message type"))
			return
		}

		callback(orderbook, nil)
	})
}
