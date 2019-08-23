package ecwid

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID = 666
		token   = "token"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/?token=%s", storeID, token)

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, req.URL.String(), expectedEndpoint, "endpoint")
			return httpmock.NewStringResponse(200, ""), nil
		})

	New(storeID, token).R().Get("/")
}
