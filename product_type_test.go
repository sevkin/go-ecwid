package ecwid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestProductTypesGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID = 666
		token   = "token"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/classes", storeID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, `[{"id":1,"name":"aa"},{"id":2,"name":"bb"}]`), nil
		})

	result, err := New(storeID, token).ProductTypesGet()
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, 2, len(*result))
	assert.Equal(t, "aa", (*result)[0].Name)
	assert.Equal(t, "bb", (*result)[1].Name)
}

func TestProductTypeGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID        = 666
		productClassID = 999
		token          = "token"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/classes/%d", storeID, productClassID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, fmt.Sprintf(`{"id":%d, "name":"name"}`, productClassID)), nil
		})

	result, err := New(storeID, token).ProductTypeGet(productClassID)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint64(productClassID), result.ID, "id")
	assert.Equal(t, "name", result.Name, "name")
}

func TestProductTypeAdd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID = 666
		token   = "token"
		name    = "new cls"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/classes", storeID)
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
			var request ProductType
			err = json.Unmarshal(body, &request)
			assert.Nil(t, err)
			assert.Equal(t, name, request.Name, "name")

			return httpmock.NewStringResponse(200, `{"id":999}`), nil
		})

	id, err := New(storeID, token).ProductTypeAdd(&ProductType{Name: name})
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint64(999), id, "id")
}

func TestProductTypeUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID        = 666
		token          = "token"
		productClassID = 999
		name           = "upd cat"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/classes/%d", storeID, productClassID)
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
			var request ProductType
			err = json.Unmarshal(body, &request)
			assert.Nil(t, err)
			assert.Equal(t, name, request.Name, "name")

			return httpmock.NewStringResponse(200, `{"updateCount":1}`), nil
		})

	err := New(storeID, token).ProductTypeUpdate(productClassID, &ProductType{Name: name})
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
}

func TestProductTypeDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID        = 666
		token          = "token"
		productClassID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/classes/%d", storeID, productClassID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "DELETE", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, `{"deleteCount":1}`), nil
		})

	err := New(storeID, token).ProductTypeDelete(productClassID)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
}
