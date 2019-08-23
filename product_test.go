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
	"github.com/stretchr/testify/assert"
)

func TestProductsSearchRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID = 666
		token   = "token"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products", storeID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")

			values := req.URL.Query()
			assert.Equal(t, "test product", values.Get("keyword"), "keyword")
			assert.Equal(t, "5", values.Get("limit"), "limit")

			return httpmock.NewStringResponse(200, ""), nil
		})

	New(storeID, token).ProductsSearch(map[string]string{
		"keyword": "test product",
		"limit":   "5",
	})
	assert.Truef(t, requested, "request failed")

}

// TODO TestProductsSearchResponse

func TestProducts(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID = 666
		token   = "token"
	)

	requestCount := 0

	expected := []string{"one", "two", "tree"}

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			if requestCount < len(expected) {

				requestCount++

				values := req.URL.Query()
				offset, _ := strconv.ParseUint(values.Get("offset"), 10, 64)
				limit, _ := strconv.ParseUint(values.Get("limit"), 10, 64)

				assert.Equal(t, uint64(requestCount), offset)
				assert.Equal(t, uint64(1), limit)

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

	for product := range New(storeID, token).Products(context.Background(), // ctx,
		filter) {
		actual = append(actual, product.Name)
	}
	assert.Equal(t, len(expected)-1, requestCount)
	assert.Equal(t, expected[1:], actual)

	assert.Equal(t, "1", filter["offset"], "filter map must be unchanged")
}

func TestProductGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID   = 666
		productID = 999
		token     = "token"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, fmt.Sprintf(`{"id":%d, "sku":"sky"}`, productID)), nil
		})

	p, err := New(storeID, token).ProductGet(productID)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint64(productID), p.ID, "id")
	assert.Equal(t, "sky", p.Sku, "sku")
}

func TestProductAdd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID = 666
		token   = "token"
		sku     = "test"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products", storeID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "POST", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")
			assert.Equal(t, "application/json", req.Header["Content-Type"][0], "Content-Type: application/json")

			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			var p Product
			err = json.Unmarshal(body, &p)
			assert.Nil(t, err)
			assert.Equal(t, sku, p.Sku, "sku")

			return httpmock.NewStringResponse(200, `{"id":999}`), nil
		})

	id, err := New(storeID, token).ProductAdd(&NewProduct{Sku: sku})
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint64(999), id, "id")
}

func TestProductUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID   = 666
		token     = "token"
		productID = 999
		sku       = "test"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "PUT", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")
			assert.Equal(t, "application/json", req.Header["Content-Type"][0], "Content-Type: application/json")

			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			var p Product
			err = json.Unmarshal(body, &p)
			assert.Nil(t, err)
			assert.Equal(t, sku, p.Sku, "sku")

			return httpmock.NewStringResponse(200, `{"updateCount":1}`), nil
		})

	err := New(storeID, token).ProductUpdate(productID, &NewProduct{Sku: sku})
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
}

func TestProductDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID   = 666
		token     = "token"
		productID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "DELETE", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, `{"deleteCount":1}`), nil
		})

	err := New(storeID, token).ProductDelete(productID)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
}

func TestProductInventoryAdjust(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID   = 666
		token     = "token"
		productID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/inventory", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "PUT", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")
			assert.Equal(t, "application/json", req.Header["Content-Type"][0], "Content-Type: application/json")

			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			var d struct {
				Delta int `json:"quantityDelta"`
			}
			err = json.Unmarshal(body, &d)
			assert.Nil(t, err)
			assert.Equal(t, -1, d.Delta, "delta")

			return httpmock.NewStringResponse(200, `{"updateCount":1}`), nil
		})

	quantity, err := New(storeID, token).ProductInventoryAdjust(productID, -1)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, 1, quantity, "quantity")
}
