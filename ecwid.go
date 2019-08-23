package ecwid

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	endpoint = "https://app.ecwid.com/api/v3/%d" // storeID
)

type (
	// Client is Ecwid API client
	Client struct {
		*resty.Client
	}
)

// New method creates a new ecwid client based on resty.Client
func New(storeID uint64, token string) *Client {
	client := resty.New().SetHostURL(fmt.Sprintf(endpoint, storeID)).SetQueryParam("token", token)
	return &Client{
		client,
	}
}
