package ecwid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestSearchProductsRequest(t *testing.T) {
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

	New(storeID, token).SearchProducts(map[string]string{
		"keyword": "test product",
		"limit":   "5",
	})
	assert.Truef(t, requested, "request failed")

}

// TODO TestSearchProductsResponse

func TestGetProduct(t *testing.T) {
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

			return httpmock.NewStringResponse(200, `{"id":999, "sku":"sky"}`), nil
		})

	p, err := New(storeID, token).GetProduct(productID)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint(999), p.ID, "id")
	assert.Equal(t, "sky", p.Sku, "sku")
}

// TODO TestGetProductResponseFail

func TestAddProduct(t *testing.T) {
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

	id, err := New(storeID, token).AddProduct(&Product{Sku: sku})
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint(999), id, "id")
}

// TODO TestAddProductResponseFail

func TestUpdateProduct(t *testing.T) {
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

	err := New(storeID, token).UpdateProduct(productID, &Product{Sku: sku})
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
}

func TestDeleteProduct(t *testing.T) {
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

	err := New(storeID, token).DeleteProduct(productID)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
}

func TestAdjustProductInventory(t *testing.T) {
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

	quantity, err := New(storeID, token).AdjustProductInventory(productID, -1)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, 1, quantity, "quantity")
}
