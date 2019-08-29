package ecwid

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// ProductImageUpload uploads product image from stream
func (c *Client) ProductImageUpload(productID ID, image io.Reader) (ID, error) {
	response, err := c.R().
		SetHeader("Content-Type", "image/jpeg").
		SetBody(image).
		Post(fmt.Sprintf("/products/%d/image", productID))

	return responseAdd(response, err)
}

// ProductImageUploadFile uploads product image from local image file
func (c *Client) ProductImageUploadFile(productID ID, filename string) (ID, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return c.ProductImageUpload(productID, bufio.NewReader(file))
}

// ProductImageUploadByURL uploads product image  from external resource
func (c *Client) ProductImageUploadByURL(productID ID, imageURL string) (ID, error) {
	response, err := c.R().
		SetQueryParam("externalUrl", imageURL).
		Post(fmt.Sprintf("/products/%d/image", productID))

	return responseAdd(response, err)
}

// ProductImageDelete deletes the main image of a product in an Ecwid store
func (c *Client) ProductImageDelete(productID ID) error {
	response, err := c.R().
		Delete(fmt.Sprintf("/products/%d/image", productID))

	_, err = responseDelete(response, err)
	return err
}
