package ecwid

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type ProductImageTestSuite struct {
	ClientTestSuite
}

func TestProductImageTestSuite(t *testing.T) {
	suite.Run(t, new(ProductImageTestSuite))
}

func (suite *ProductImageTestSuite) TestProductImageUpload() {
	const (
		productID = 999
		imageFile = "fixture/ecwid.jpg"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/image", storeID, productID)
	requested := false

	file, err := os.Open(imageFile)
	suite.Nil(err)
	defer file.Close()
	image, err := ioutil.ReadAll(file)
	suite.Nil(err)
	file.Seek(0, 0)

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("POST", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")
			suite.Equal("image/jpeg", req.Header["Content-Type"][0], "Content-Type: image/jpeg")

			body, err := ioutil.ReadAll(req.Body)
			suite.Nil(err)
			suite.Equal(image, body)

			return httpmock.NewStringResponse(200, `{"id":999}`), nil
		})

	id, err := suite.client.ProductImageUpload(productID, file)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(uint64(999), id, "id")
}

func (suite *ProductImageTestSuite) TestProductImageUploadFile() {
	const (
		productID = 999
		imageFile = "fixture/ecwid.jpg"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/image", storeID, productID)
	requested := false

	file, err := os.Open(imageFile)
	suite.Nil(err)
	defer file.Close()
	image, err := ioutil.ReadAll(file)
	suite.Nil(err)

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("POST", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")
			suite.Equal("image/jpeg", req.Header["Content-Type"][0], "Content-Type: image/jpeg")

			body, err := ioutil.ReadAll(req.Body)
			suite.Nil(err)
			suite.Equal(image, body)

			return httpmock.NewStringResponse(200, `{"id":999}`), nil
		})

	id, err := suite.client.ProductImageUploadFile(productID, imageFile)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(uint64(999), id, "id")
}

func (suite *ProductImageTestSuite) TestProductImageUploadFileNotFound() {
	const (
		productID = 999
		imageFile = "fixture/notfound.jpg"
	)

	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			return httpmock.NewStringResponse(400, `{"errorMessage":"ignore me"}`), nil
		})

	_, err := suite.client.ProductImageUploadFile(productID, imageFile)
	suite.NotNil(err)
	suite.Falsef(requested, "request failed")
}

func (suite *ProductImageTestSuite) TestProductImageUploadByURL() {
	const (
		productID = 999
		imageURL  = "https://example.org/image.jpg"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/image", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("POST", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")
			values := req.URL.Query()

			suite.Equal(imageURL, values.Get("externalUrl"), "externalUrl")

			return httpmock.NewStringResponse(200, `{"id":999}`), nil
		})

	id, err := suite.client.ProductImageUploadByURL(productID, imageURL)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(uint64(999), id, "id")
}

func (suite *ProductImageTestSuite) TestProductImageDelete() {
	const (
		productID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/image", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("DELETE", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")
			// suite.Equal("application/json", req.Header["Content-Type"][0], "Content-Type: application/json")

			return httpmock.NewStringResponse(200, `{"deleteCount":1}`), nil
		})

	err := suite.client.ProductImageDelete(productID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
}
