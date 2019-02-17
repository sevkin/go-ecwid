package ecwid

import (
	"fmt"

	resty "gopkg.in/resty.v1"
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
func New(storeID uint, token string) *Client {
	client := resty.New().SetHostURL(fmt.Sprintf(endpoint, storeID)).SetQueryParam("token", token)
	return &Client{
		client,
	}
}
