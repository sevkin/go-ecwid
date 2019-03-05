package ecwid

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// ProductImageGalleryUpload uploads image to product gallery from stream
func (c *Client) ProductImageGalleryUpload(productID uint64, image io.Reader, imageTitle string) (uint64, error) {
	params := make(map[string]string)
	if len(imageTitle) > 0 {
		params["fileName"] = imageTitle
	}

	response, err := c.R().
		SetQueryParams(params).
		SetHeader("Content-Type", "image/jpeg").
		SetBody(image).
		Post(fmt.Sprintf("/products/%d/gallery", productID))

	return responseAdd(response, err)
}

// ProductImageGalleryUploadFile uploads image to product gallery from local file
func (c *Client) ProductImageGalleryUploadFile(productID uint64, filename string, imageTitle string) (uint64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return c.ProductImageGalleryUpload(productID, bufio.NewReader(file), imageTitle)
}

// ProductImageGalleryUploadByURL uploads image to product gallery from external url
func (c *Client) ProductImageGalleryUploadByURL(productID uint64, imageURL, imageTitle string) (uint64, error) {
	params := map[string]string{
		"externalUrl": imageURL,
	}
	if len(imageTitle) > 0 {
		params["fileName"] = imageTitle
	}

	response, err := c.R().
		SetQueryParams(params).
		Post(fmt.Sprintf("/products/%d/gallery", productID))

	return responseAdd(response, err)
}

// ProductImageGalleryDelete deletes image from product gallery by image id
func (c *Client) ProductImageGalleryDelete(productID, fieldID uint64) error {
	response, err := c.R().
		Delete(fmt.Sprintf("/products/%d/gallery/%d", productID, fieldID))

	_, err = responseDelete(response, err)
	return err
}

// ProductImageGalleryDeleteAll deletes all images from product gallery
func (c *Client) ProductImageGalleryDeleteAll(productID uint64) (uint, error) {
	response, err := c.R().
		Delete(fmt.Sprintf("/products/%d/gallery", productID))

	return responseDelete(response, err)
}
