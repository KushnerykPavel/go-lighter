package lighter

import (
	"net/http"
)

const (
	MainnetAPIURL     = "https://mainnet.zklighter.elliot.ai"
	MainnetRestAPIURL = "https://mainnet.zklighter.elliot.ai/api/v1/"
)

type Client struct {
	apiURL string
	c      *http.Client
}

func NewClient(apiURL string, httpClient *http.Client) *Client {
	u := apiURL
	if u == "" {
		u = MainnetAPIURL
	}
	return &Client{
		apiURL: u,
		c:      httpClient,
	}
}
