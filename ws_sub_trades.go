package lighter

import (
	"fmt"
)

type TradesSubscriptionParams struct {
	Coin string
}

func (w *WebsocketClient) Trades(
	params TradesSubscriptionParams,
	callback func(*Trades, error),
) (*Subscription, error) {
	remotePayload := remoteTradesSubscriptionPayload{
		Type:      "subscribe",
		Operation: ChannelTrades,
		Coin:      params.Coin,
		Channel:   fmt.Sprintf("%s/%s", ChannelTrades, params.Coin),
	}
	return w.subscribe(remotePayload, func(msg any) {
		trades, ok := msg.(*Trades)
		if !ok {
			callback(&Trades{}, fmt.Errorf("invalid message type"))
			return
		}

		callback(trades, nil)
	})
}
