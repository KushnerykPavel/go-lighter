package lighter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type debugLogger struct{}

func (d debugLogger) Infof(format string, args ...any) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func (d debugLogger) Errorf(format string, args ...any) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}

func TestWebsocketClient_L2Book(t *testing.T) {
	client, err := NewWebsocketClient("")
	assert.NoError(t, err)

	err = client.Connect(context.Background())
	assert.NoError(t, err)

	_, _ = client.OrderBook(OrderBookSubscriptionParams{Coin: "15"}, func(book *OrderBook, err error) {
		fmt.Println(book)
	})

	time.Sleep(time.Minute)
}
