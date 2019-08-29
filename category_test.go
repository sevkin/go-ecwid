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

type CategoryTestSuite struct {
	ClientTestSuite
}

func TestCategoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryTestSuite))
}

func (suite *CategoryTestSuite) TestCategoriesGet() {
	const (
		url = "http://example.org/cat"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/categories", storeID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			values := req.URL.Query()
			suite.Equal("0", values.Get("parent"), "parent")
			suite.Equal("5", values.Get("limit"), "limit")

			return httpmock.NewStringResponse(200, `{"total":2,"count":2,"offset":0,"limit":100,"items":[{"id": 1, "name": "one"},{"id": 2, "url": "`+url+`"}]}`), nil
		})

	result, err := suite.client.CategoriesGet(map[string]string{
		"parent": "0",
		"limit":  "5",
	})
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.NotNil(result)

	suite.Equal(2, len(result.Items))
	suite.Equal("one", result.Items[0].Name)
	suite.Equal(url, result.Items[1].URL)
}

func (suite *CategoryTestSuite) TestCategories() {
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

				return httpmock.NewJsonResponse(200, CategoriesGetResponse{
					Total:  3,
					Count:  1,
					Offset: uint(offset),
					Limit:  uint(limit),
					Items: []*Category{
						&Category{
							NewCategory: NewCategory{
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

	for category := range suite.client.Categories(context.Background(), // ctx,
		filter) {
		actual = append(actual, category.Name)
	}
	suite.Equal(len(expected)-1, requestCount)
	suite.Equal(expected[1:], actual)

	suite.Equal("1", filter["offset"], "filter map must be unchanged")
}

func (suite *CategoryTestSuite) TestCategoryGet() {
	const (
		categoryID ID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/categories/%d", storeID, categoryID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, fmt.Sprintf(`{"id":%d, "name":"name"}`, categoryID)), nil
		})

	c, err := suite.client.CategoryGet(categoryID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(categoryID, c.ID, "id")
	suite.Equal("name", c.Name, "name")
}

func (suite *CategoryTestSuite) TestCategoryAdd() {
	const (
		categoryID ID = 999
		name          = "new cat"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/categories", storeID)
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
			var c NewCategory
			err = json.Unmarshal(body, &c)
			suite.Nil(err)
			suite.Equal(name, c.Name, "name")

			return httpmock.NewStringResponse(200, fmt.Sprintf(`{"id":%d}`, categoryID)), nil
		})

	id, err := suite.client.CategoryAdd(&NewCategory{Name: name})
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(categoryID, id, "id")
}

func (suite *CategoryTestSuite) TestCategoryUpdate() {
	const (
		categoryID ID = 999
		name          = "upd cat"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/categories/%d", storeID, categoryID)
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
			var c NewCategory
			err = json.Unmarshal(body, &c)
			suite.Nil(err)
			suite.Equal(name, c.Name, "name")

			return httpmock.NewStringResponse(200, `{"updateCount":1}`), nil
		})

	err := suite.client.CategoryUpdate(categoryID, &NewCategory{Name: name})
	suite.Truef(requested, "request failed")

	suite.Nil(err)
}

func (suite *CategoryTestSuite) TestCategoryDelete() {
	const (
		categoryID ID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/categories/%d", storeID, categoryID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("DELETE", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, `{"deleteCount":1}`), nil
		})

	err := suite.client.CategoryDelete(categoryID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
}
