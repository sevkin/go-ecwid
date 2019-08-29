package ecwid

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type StoreInformationTestSuite struct {
	ClientTestSuite
}

func TestStoreInformationTestSuite(t *testing.T) {
	suite.Run(t, new(StoreInformationTestSuite))
}

func (suite *StoreInformationTestSuite) TestStoreProfileGet() {
	expectedEndpoint := fmt.Sprintf(endpoint+"/profile", storeID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200,
				fmt.Sprintf(`{"generalInfo": {"storeId": %d}}`, storeID)), nil
		})

	p, err := suite.client.StoreProfileGet()
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(storeID, p.GeneralInfo.StoreID, "id")
}
