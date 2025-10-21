package lighter

import (
	"strings"
)

func key(args ...string) string {
	return strings.Join(args, ":")
}

func keyOrderBook(coin string) string {
	return key(ChannelOrderBook, coin)
}

func keyTrades(coin string) string {
	return key(ChannelTrades, coin)
}
