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

type ProductTypeTestSuite struct {
	ClientTestSuite
}

func TestProductTypeTestSuite(t *testing.T) {
	suite.Run(t, new(ProductTypeTestSuite))
}

func (suite *ProductTypeTestSuite) TestProductTypesGet() {
	expectedEndpoint := fmt.Sprintf(endpoint+"/classes", storeID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, `[{"id":1,"name":"aa"},{"id":2,"name":"bb"}]`), nil
		})

	result, err := suite.client.ProductTypesGet()
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.NotNil(result)

	suite.Equal(2, len(*result))
	suite.Equal("aa", (*result)[0].Name)
	suite.Equal("bb", (*result)[1].Name)
}

func (suite *ProductTypeTestSuite) TestProductTypeGet() {
	const (
		productClassID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/classes/%d", storeID, productClassID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, fmt.Sprintf(`{"id":%d, "name":"name"}`, productClassID)), nil
		})

	result, err := suite.client.ProductTypeGet(productClassID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(uint64(productClassID), result.ID, "id")
	suite.Equal("name", result.Name, "name")
}

func (suite *ProductTypeTestSuite) TestProductTypeAdd() {
	const (
		name = "new cls"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/classes", storeID)
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
			var request ProductType
			err = json.Unmarshal(body, &request)
			suite.Nil(err)
			suite.Equal(name, request.Name, "name")

			return httpmock.NewStringResponse(200, `{"id":999}`), nil
		})

	id, err := suite.client.ProductTypeAdd(&ProductType{Name: name})
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(uint64(999), id, "id")
}

func (suite *ProductTypeTestSuite) TestProductTypeUpdate() {
	const (
		productClassID = 999
		name           = "upd cat"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/classes/%d", storeID, productClassID)
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
			var request ProductType
			err = json.Unmarshal(body, &request)
			suite.Nil(err)
			suite.Equal(name, request.Name, "name")

			return httpmock.NewStringResponse(200, `{"updateCount":1}`), nil
		})

	err := suite.client.ProductTypeUpdate(productClassID, &ProductType{Name: name})
	suite.Truef(requested, "request failed")

	suite.Nil(err)
}

func (suite *ProductTypeTestSuite) TestProductTypeDelete() {
	const (
		productClassID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/classes/%d", storeID, productClassID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("DELETE", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200, `{"deleteCount":1}`), nil
		})

	err := suite.client.ProductTypeDelete(productClassID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
}
