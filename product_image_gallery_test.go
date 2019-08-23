package ecwid

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestProductImageGalleryUpload(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID    = 666
		token      = "token"
		productID  = 999
		imageFile  = "fixture/ecwid.jpg"
		imageTitle = "test gallery upload"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/gallery", storeID, productID)
	requested := false

	file, err := os.Open(imageFile)
	assert.Nil(t, err)
	defer file.Close()
	image, err := ioutil.ReadAll(file)
	assert.Nil(t, err)
	file.Seek(0, 0)

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "POST", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")
			assert.Equal(t, "image/jpeg", req.Header["Content-Type"][0], "Content-Type: image/jpeg")

			values := req.URL.Query()
			assert.Equal(t, imageTitle, values.Get("fileName"), "fileName")

			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.Equal(t, image, body)

			return httpmock.NewStringResponse(200, `{"id":999}`), nil
		})

	id, err := New(storeID, token).ProductImageGalleryUpload(productID, file, imageTitle)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint64(999), id, "id")
}

func TestProductImageGalleryUploadFile(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID   = 666
		token     = "token"
		productID = 999
		imageFile = "fixture/ecwid.jpg"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/gallery", storeID, productID)
	requested := false

	file, err := os.Open(imageFile)
	assert.Nil(t, err)
	defer file.Close()
	image, err := ioutil.ReadAll(file)
	assert.Nil(t, err)

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "POST", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")
			assert.Equal(t, "image/jpeg", req.Header["Content-Type"][0], "Content-Type: image/jpeg")

			values := req.URL.Query()
			_, ok := values["fileName"]
			assert.Falsef(t, ok, "fileName")

			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.Equal(t, image, body)

			return httpmock.NewStringResponse(200, `{"id":999}`), nil
		})

	id, err := New(storeID, token).ProductImageGalleryUploadFile(productID, imageFile, "")
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint64(999), id, "id")
}

func TestProductImageGalleryUploadFileNotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID   = 666
		token     = "token"
		productID = 999
		imageFile = "fixture/notfound.jpg"
	)

	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			return httpmock.NewStringResponse(400, `{"errorMessage":"ignore me"}`), nil
		})

	_, err := New(storeID, token).ProductImageGalleryUploadFile(productID, imageFile, "")
	assert.NotNil(t, err)
	assert.Falsef(t, requested, "request failed")
}

func TestProductImageGalleryUploadByURL(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID   = 666
		token     = "token"
		productID = 999
		imageURL  = "https://example.org/image.jpg"
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/gallery", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "POST", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")
			values := req.URL.Query()

			assert.Equal(t, imageURL, values.Get("externalUrl"), "externalUrl")

			return httpmock.NewStringResponse(200, `{"id":999}`), nil
		})

	id, err := New(storeID, token).ProductImageGalleryUploadByURL(productID, imageURL, "")
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint64(999), id, "id")
}

func TestProductImageGalleryDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID   = 666
		token     = "token"
		productID = 999
		fieldID   = 555
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/gallery/%d", storeID, productID, fieldID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "DELETE", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")
			// assert.Equal(t, "application/json", req.Header["Content-Type"][0], "Content-Type: application/json")

			return httpmock.NewStringResponse(200, `{"deleteCount":1}`), nil
		})

	err := New(storeID, token).ProductImageGalleryDelete(productID, fieldID)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
}

func TestProductImageGalleryDeleteAll(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		storeID   = 666
		token     = "token"
		productID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/products/%d/gallery", storeID, productID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			assert.Equal(t, "DELETE", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			assert.Equal(t, expectedEndpoint, actualEndpoint, "endpoint")
			// assert.Equal(t, "application/json", req.Header["Content-Type"][0], "Content-Type: application/json")

			return httpmock.NewStringResponse(200, `{"deleteCount":4}`), nil
		})

	cnt, err := New(storeID, token).ProductImageGalleryDeleteAll(productID)
	assert.Truef(t, requested, "request failed")

	assert.Nil(t, err)
	assert.Equal(t, uint(4), cnt, "deleteCount")
}
