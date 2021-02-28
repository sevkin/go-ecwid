package ecwid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type ProductVariationTestSuite struct {
	ClientTestSuite
}

func TestProductVariationTestSuite(t *testing.T) {
	suite.Run(t, new(ProductVariationTestSuite))
}

func (suite *ProductVariationTestSuite) TestProductVariationsGet() {
	const (
		productID    ID = 999
		variationID1 ID = 555
		variationID2 ID = 777
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/combinations", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, fmt.Sprintf(`[{"id":%d},{"id":%d}]`, variationID1, variationID2)), nil
		})

	pv, err := suite.client.ProductVariationsGet(productID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(variationID1, pv[0].ID, "id 1")
	suite.Equal(variationID2, pv[1].ID, "id 2")
}

func (suite *ProductVariationTestSuite) TestProductVariationGet() {
	const (
		productID   ID = 999
		variationID ID = 555
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/combinations/%d", storeID, productID, variationID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, fmt.Sprintf(`{"id":%d}`, variationID)), nil
		})

	pv, err := suite.client.ProductVariationGet(productID, variationID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(variationID, pv.ID, "id")
}

func (suite *ProductVariationTestSuite) TestProductVariationUpdate() {
	const (
		productID   ID = 999
		variationID ID = 555
		sku            = "test"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/combinations/%d", storeID, productID, variationID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("PUT", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")
			suite.Equal("application/json", req.Header["Content-Type"][0], "Content-Type: application/json")

			body, err := ioutil.ReadAll(req.Body)
			suite.Nil(err)
			var pv NewProductVariation
			err = json.Unmarshal(body, &pv)
			suite.Nil(err)
			suite.Equal(sku, pv.Sku, "sku")

			return httpmock.NewStringResponse(200, `{"updateCount":1}`), nil
		})

	err := suite.client.ProductVariationUpdate(productID, variationID, &NewProductVariation{Sku: sku})
	suite.Truef(requested, "request failed")

	suite.Nil(err)
}
