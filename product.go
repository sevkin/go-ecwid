package ecwid

import (
	"fmt"
)

type (
	// Product https://developers.ecwid.com/api-documentation/products
	Product struct {
		ID                             uint    `json:"id"`
		Sku                            string  `json:"sku"`
		Quantity                       int     `json:"quantity"`
		Unlimited                      bool    `json:"unlimited"`
		InStock                        bool    `json:"inStock"`
		Name                           string  `json:"name"`
		Price                          float32 `json:"price"`
		DefaultDisplayedPrice          float32 `json:"defaultDisplayedPrice"`
		DefaultDisplayedPriceFormatted string  `json:"defaultDisplayedPriceFormatted"`
		// TODO: add more filelds
	}

	// SearchProductsResponse https://developers.ecwid.com/api-documentation/products#search-products
	SearchProductsResponse struct {
		Total    uint       `json:"total"`
		Count    uint       `json:"count"`
		Offset   uint       `json:"offset"`
		Limit    uint       `json:"limit"`
		Products []*Product `json:"items"`
	}
)

// SearchProducts search or filter products in a store catalog
func (c *Client) SearchProducts(filter map[string]string) (*SearchProductsResponse, error) {
	// filter:
	// keyword string, priceFrom number, priceTo number, category number,
	// withSubcategories bool, sortBy enum, offset number, limit number,
	// createdFrom date, createdTo date, updatedFrom date, updatedTo date,
	// enabled bool, inStock bool, onsale string,
	// sku string, productId number, baseUrl string, cleanUrls bool,
	// TODO Search Products: field, option, attribute ???
	// field{attributeName}={attributeValues} field{attributeId}={attributeValues}
	// option_{optionName}={optionValues}
	// attribute_{attributeName}={attributeValues}

	response, err := c.R().
		SetQueryParams(filter).
		Get("/products")

	var result SearchProductsResponse
	return &result, responseUnmarshal(response, err, &result)
}

// GetProduct gets all details of a specific product in an Ecwid store by its ID
func (c *Client) GetProduct(productID uint) (*Product, error) {
	response, err := c.R().
		SetPathParams(map[string]string{
			"productId": fmt.Sprintf("%d", productID),
		}).
		Get("/products/{productId}")

	var result Product
	return &result, responseUnmarshal(response, err, &result)
}

// AddProduct creates a new product in an Ecwid store
// returns new productId
func (c *Client) AddProduct(product *Product) (uint, error) {
	// FIXME!!! похоже на каждый запрос (add|get|search|update) набор полей различный

	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(product).
		Post("/products")

	id, err := responseAdd(response, err)
	return uint(id), err // TODO ID uint64
}

// UpdateProduct update an existing product in an Ecwid store referring to its ID
// before update use GetProduct to retrieve full data
func (c *Client) UpdateProduct(productID uint, product *Product) error {
	response, err := c.R().
		SetPathParams(map[string]string{
			"productId": fmt.Sprintf("%d", productID),
		}).
		SetHeader("Content-Type", "application/json").
		SetBody(product).
		Put("/products/{productId}")

	return responseUpdate(response, err)
}

// TODO try to pass partial json with help https://github.com/tidwall/sjson
// func (c *Client) UpdateProductJson(productID uint, productJson string) error {

// DeleteProduct delete a product from an Ecwid store referring to its ID
func (c *Client) DeleteProduct(productID uint) error {
	response, err := c.R().
		SetPathParams(map[string]string{
			"productId": fmt.Sprintf("%d", productID),
		}).
		Delete("/products/{productId}")

	return responseDelete(response, err)
}

// AdjustProductInventory increase or decrease the product’s stock quantity by a delta quantity
func (c *Client) AdjustProductInventory(productID uint, quantityDelta int) (int, error) {
	response, err := c.R().
		SetPathParams(map[string]string{
			"productId": fmt.Sprintf("%d", productID),
		}).
		SetBody(fmt.Sprintf(`{"quantityDelta":%d}`, quantityDelta)).
		Put("/products/{productId}/inventory")

	return responseUpdateCount(response, err)
}
