package ecwid

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type ProductTestSuite struct {
	ClientTestSuite
}

func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}

func (suite *ProductTestSuite) TestProductsSearchRequest() {
	expectedEndpoint := fmt.Sprintf(endpoint+"/products", storeID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			values := req.URL.Query()
			suite.Equal("test product", values.Get("keyword"), "keyword")
			suite.Equal("5", values.Get("limit"), "limit")

			return httpmock.NewStringResponse(200, ""), nil
		})

	suite.client.ProductsSearch(map[string]string{
		"keyword": "test product",
		"limit":   "5",
	})
	suite.Truef(requested, "request failed")

}

// TODO TestProductsSearchResponse

func (suite *ProductTestSuite) TestProducts() {
	requestCount := 0

	expected := []string{"one", "two", "tree"}

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			if requestCount < len(expected) {

				requestCount++

				values := req.URL.Query()
				offset, _ := strconv.ParseUint(values.Get("offset"), 10, 64)
				limit, _ := strconv.ParseUint(values.Get("limit"), 10, 64)

				suite.Equal(uint64(requestCount), offset)
				suite.Equal(uint64(1), limit)

				return httpmock.NewJsonResponse(200, ProductsSearchResponse{
					Total:  3,
					Count:  1,
					Offset: uint(offset),
					Limit:  uint(limit),
					Products: []*Product{
						&Product{
							NewProduct: NewProduct{
								Name: expected[offset],
							},
						},
					},
				})
			}
			return httpmock.NewStringResponse(400, "too many"), nil
		})

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	actual := make([]string, 0, 3)

	filter := map[string]string{
		"offset": "1",
		"limit":  "1",
	}

	for product := range suite.client.Products(context.Background(), // ctx,
		filter) {
		actual = append(actual, product.Name)
	}
	suite.Equal(len(expected)-1, requestCount)
	suite.Equal(expected[1:], actual)

	suite.Equal("1", filter["offset"], "filter map must be unchanged")
}

func (suite *ProductTestSuite) TestProductGet() {
	const (
		productID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, fmt.Sprintf(`{"id":%d, "sku":"sky"}`, productID)), nil
		})

	p, err := suite.client.ProductGet(productID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(uint64(productID), p.ID, "id")
	suite.Equal("sky", p.Sku, "sku")
}

func (suite *ProductTestSuite) TestProductAdd() {
	const (
		sku = "test"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products", storeID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("POST", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")
			suite.Equal("application/json", req.Header["Content-Type"][0], "Content-Type: application/json")

			body, err := ioutil.ReadAll(req.Body)
			suite.Nil(err)
			var p Product
			err = json.Unmarshal(body, &p)
			suite.Nil(err)
			suite.Equal(sku, p.Sku, "sku")

			return httpmock.NewStringResponse(200, `{"id":999}`), nil
		})

	id, err := suite.client.ProductAdd(&NewProduct{Sku: sku})
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(uint64(999), id, "id")
}

func (suite *ProductTestSuite) TestProductUpdate() {
	const (
		productID = 999
		sku       = "test"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d", storeID, productID)
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
			var p Product
			err = json.Unmarshal(body, &p)
			suite.Nil(err)
			suite.Equal(sku, p.Sku, "sku")

			return httpmock.NewStringResponse(200, `{"updateCount":1}`), nil
		})

	err := suite.client.ProductUpdate(productID, &NewProduct{Sku: sku})
	suite.Truef(requested, "request failed")

	suite.Nil(err)
}

func (suite *ProductTestSuite) TestProductDelete() {
	const (
		productID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("DELETE", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, `{"deleteCount":1}`), nil
		})

	err := suite.client.ProductDelete(productID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
}

func (suite *ProductTestSuite) TestProductInventoryAdjust() {
	const (
		productID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/inventory", storeID, productID)
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
			var d struct {
				Delta int `json:"quantityDelta"`
			}
			err = json.Unmarshal(body, &d)
			suite.Nil(err)
			suite.Equal(-1, d.Delta, "delta")

			return httpmock.NewStringResponse(200, `{"updateCount":1}`), nil
		})

	quantity, err := suite.client.ProductInventoryAdjust(productID, -1)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(1, quantity, "quantity")
}
