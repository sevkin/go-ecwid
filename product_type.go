package ecwid

// Product types (or product classes) are groups of products which share the same attributes.

import "fmt"

type (
	// ProductType https://developers.ecwid.com/api-documentation/product-types#get-product-type
	ProductType struct {
		ID         ID         `json:"id,omitempty"` // mandatory for update
		Name       string     `json:"name,omitempty"`
		Attributes Attributes `json:"attributes,omitempty"`
	}

	// ProductTypesResponse https://developers.ecwid.com/api-documentation/product-types#get-product-types
	ProductTypesResponse []ProductType
)

// ProductTypesGet gets all product types present in an Ecwid store
func (c *Client) ProductTypesGet() (*ProductTypesResponse, error) {
	response, err := c.R().
		Get("/classes")

	var result ProductTypesResponse
	return &result, responseUnmarshal(response, err, &result)
}

// ProductTypeGet gets the full details of a specific product type referring to its ID
func (c *Client) ProductTypeGet(productClassID ID) (*ProductType, error) {
	response, err := c.R().
		Get(fmt.Sprintf("/classes/%d", productClassID))

	var result ProductType
	return &result, responseUnmarshal(response, err, &result)
}

// ProductTypeAdd creates a new product type in an Ecwid store
// returns new productClassID
func (c *Client) ProductTypeAdd(ProductType *ProductType) (ID, error) {
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(ProductType).
		Post("/classes")

	return responseAdd(response, err)
}

// ProductTypeUpdate updates the details of a specific product type referring to its ID.
// If you need to update existing product attributes, refer to their IDs in your request.
// If you want to add new product attributes to existing product type,
// send your new attributes AND all existing attributes for that product type.
// Otherwise you will reset the existing attributes in that product type.
func (c *Client) ProductTypeUpdate(productClassID ID, ProductType *ProductType) error {
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(ProductType).
		Put(fmt.Sprintf("/classes/%d", productClassID))

	return responseUpdate(response, err)
}

// ProductTypeDelete deletes a specific product type and its assigned attributes.
// The products that belong to this type will not be removed.
// They will be re-assigned to the General type.
func (c *Client) ProductTypeDelete(productClassID ID) error {
	response, err := c.R().
		Delete(fmt.Sprintf("/classes/%d", productClassID))

	_, err = responseDelete(response, err)
	return err
}
