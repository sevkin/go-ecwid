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

type ProductImageGalleryTestSuite struct {
	ClientTestSuite
}

func TestProductImageGalleryTestSuite(t *testing.T) {
	suite.Run(t, new(ProductImageGalleryTestSuite))
}

func (suite *ProductImageGalleryTestSuite) TestProductImageGalleryUpload() {
	const (
		productID  ID = 999
		imageFile     = "fixture/ecwid.jpg"
		imageTitle    = "test gallery upload"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/gallery", storeID, productID)
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

			values := req.URL.Query()
			suite.Equal(imageTitle, values.Get("fileName"), "fileName")

			body, err := ioutil.ReadAll(req.Body)
			suite.Nil(err)
			suite.Equal(image, body)

			return httpmock.NewStringResponse(200, fmt.Sprintf(`{"id":%d}`, productID)), nil
		})

	id, err := suite.client.ProductImageGalleryUpload(productID, file, imageTitle)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(productID, id, "id")
}

func (suite *ProductImageGalleryTestSuite) TestProductImageGalleryUploadFile() {
	const (
		productID ID = 999
		imageFile    = "fixture/ecwid.jpg"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/gallery", storeID, productID)
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

			values := req.URL.Query()
			_, ok := values["fileName"]
			suite.Falsef(ok, "fileName")

			body, err := ioutil.ReadAll(req.Body)
			suite.Nil(err)
			suite.Equal(image, body)

			return httpmock.NewStringResponse(200, fmt.Sprintf(`{"id":%d}`, productID)), nil
		})

	id, err := suite.client.ProductImageGalleryUploadFile(productID, imageFile, "")
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(productID, id, "id")
}

func (suite *ProductImageGalleryTestSuite) TestProductImageGalleryUploadFileNotFound() {
	const (
		productID ID = 999
		imageFile    = "fixture/notfound.jpg"
	)

	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			return httpmock.NewStringResponse(400, `{"errorMessage":"ignore me"}`), nil
		})

	_, err := suite.client.ProductImageGalleryUploadFile(productID, imageFile, "")
	suite.NotNil(err)
	suite.Falsef(requested, "request failed")
}

func (suite *ProductImageGalleryTestSuite) TestProductImageGalleryUploadByURL() {
	const (
		productID ID = 999
		imageURL     = "https://example.org/image.jpg"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/gallery", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("POST", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")
			values := req.URL.Query()

			suite.Equal(imageURL, values.Get("externalUrl"), "externalUrl")

			return httpmock.NewStringResponse(200, fmt.Sprintf(`{"id":%d}`, productID)), nil
		})

	id, err := suite.client.ProductImageGalleryUploadByURL(productID, imageURL, "")
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(productID, id, "id")
}

func (suite *ProductImageGalleryTestSuite) TestProductImageGalleryDelete() {
	const (
		productID ID = 999
		fieldID   ID = 555
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/gallery/%d", storeID, productID, fieldID)
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

	err := suite.client.ProductImageGalleryDelete(productID, fieldID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
}

func (suite *ProductImageGalleryTestSuite) TestProductImageGalleryDeleteAll() {
	const (
		productID ID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/gallery", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("DELETE", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")
			// suite.Equal("application/json", req.Header["Content-Type"][0], "Content-Type: application/json")

			return httpmock.NewStringResponse(200, `{"deleteCount":4}`), nil
		})

	cnt, err := suite.client.ProductImageGalleryDeleteAll(productID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(uint(4), cnt, "deleteCount")
}
