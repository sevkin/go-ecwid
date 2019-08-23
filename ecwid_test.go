package ecwid

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type EcwidTestSuite struct {
	ClientTestSuite
}

func TestEcwidTestSuite(t *testing.T) {
	suite.Run(t, new(EcwidTestSuite))
}

func (suite *EcwidTestSuite) TestRequest() {
	expectedEndpoint := fmt.Sprintf(endpoint+"/?token=%s", storeID, token)

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			suite.Equal(req.URL.String(), expectedEndpoint, "endpoint")
			return httpmock.NewStringResponse(200, ""), nil
		})

	suite.client.R().Get("/")
}
