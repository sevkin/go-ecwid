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

func TestCategoriesGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID = 666
		token   = "token"
		url     = "http://example.org/cat"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/categories", storeID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")

			values := req.URL.Query()
			assert.Equal(t, "0", values.Get("parent"), "parent")
			assert.Equal(t, "5", values.Get("limit"), "limit")

			return httpmock.NewStringResponse(200, `{"total":2,"count":2,"offset":0,"limit":100,"items":[{"id": 1, "name": "one"},{"id": 2, "url": "`+url+`"}]}`), nil
		})

	result, err := New(storeID, token).CategoriesGet(map[string]string{
		"parent": "0",
		"limit":  "5",
	})
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, 2, len(result.Items))
	assert.Equal(t, "one", result.Items[0].Name)
	assert.Equal(t, url, result.Items[1].URL)
}

func TestCategoryGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID    = 666
		categoryID = 999
		token      = "token"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/categories/%d", storeID, categoryID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, fmt.Sprintf(`{"id":%d, "name":"name"}`, categoryID)), nil
		})

	c, err := New(storeID, token).CategoryGet(categoryID)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint64(categoryID), c.ID, "id")
	assert.Equal(t, "name", c.Name, "name")
}

func TestCategoryAdd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID = 666
		token   = "token"
		name    = "new cat"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/categories", storeID)
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
			var c NewCategory
			err = json.Unmarshal(body, &c)
			assert.Nil(t, err)
			assert.Equal(t, name, c.Name, "name")

			return httpmock.NewStringResponse(200, `{"id":999}`), nil
		})

	id, err := New(storeID, token).CategoryAdd(&NewCategory{Name: name})
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint64(999), id, "id")
}

func TestCategoryUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID    = 666
		token      = "token"
		categoryID = 999
		name       = "upd cat"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/categories/%d", storeID, categoryID)
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
			var c NewCategory
			err = json.Unmarshal(body, &c)
			assert.Nil(t, err)
			assert.Equal(t, name, c.Name, "name")

			return httpmock.NewStringResponse(200, `{"updateCount":1}`), nil
		})

	err := New(storeID, token).CategoryUpdate(categoryID, &NewCategory{Name: name})
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
}

func TestCategoryDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID    = 666
		token      = "token"
		categoryID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/categories/%d", storeID, categoryID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "DELETE", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, `{"deleteCount":1}`), nil
		})

	err := New(storeID, token).CategoryDelete(categoryID)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
}
