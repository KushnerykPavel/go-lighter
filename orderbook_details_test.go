package lighter

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_OrderBookDetails(t *testing.T) {
	client := NewClient(MainnetRestAPIURL, http.DefaultClient)
	ctx := context.Background()

	rates, err := client.OrderBookDetails(ctx)
	require.NoError(t, err)
	require.Greater(t, len(rates), 0)
}
