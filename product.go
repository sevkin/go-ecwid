package ecwid

import (
	"fmt"
)

type (
	// Product https://developers.ecwid.com/api-documentation/products
	Product struct {
		ID                             uint64  `json:"id"`
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

	// ProductsSearchResponse https://developers.ecwid.com/api-documentation/products#search-products
	ProductsSearchResponse struct {
		Total    uint       `json:"total"`
		Count    uint       `json:"count"`
		Offset   uint       `json:"offset"`
		Limit    uint       `json:"limit"`
		Products []*Product `json:"items"`
	}
)

// ProductsSearch search or filter products in a store catalog
func (c *Client) ProductsSearch(filter map[string]string) (*ProductsSearchResponse, error) {
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

	var result ProductsSearchResponse
	return &result, responseUnmarshal(response, err, &result)
}

// ProductGet gets all details of a specific product in an Ecwid store by its ID
func (c *Client) ProductGet(productID uint64) (*Product, error) {
	response, err := c.R().
		SetPathParams(map[string]string{
			"productId": fmt.Sprintf("%d", productID),
		}).
		Get("/products/{productId}")

	var result Product
	return &result, responseUnmarshal(response, err, &result)
}

// ProductAdd creates a new product in an Ecwid store
// returns new productId
func (c *Client) ProductAdd(product *Product) (uint64, error) {
	// FIXME!!! похоже на каждый запрос (add|get|search|update) набор полей различный

	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(product).
		Post("/products")

	return responseAdd(response, err)
}

// ProductUpdate update an existing product in an Ecwid store referring to its ID
// before update use ProductGet to retrieve full data
func (c *Client) ProductUpdate(productID uint64, product *Product) error {
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

// ProductDelete delete a product from an Ecwid store referring to its ID
func (c *Client) ProductDelete(productID uint64) error {
	response, err := c.R().
		SetPathParams(map[string]string{
			"productId": fmt.Sprintf("%d", productID),
		}).
		Delete("/products/{productId}")

	return responseDelete(response, err)
}

// ProductInventoryAdjust increase or decrease the product’s stock quantity by a delta quantity
func (c *Client) ProductInventoryAdjust(productID uint, quantityDelta int) (int, error) {
	response, err := c.R().
		SetPathParams(map[string]string{
			"productId": fmt.Sprintf("%d", productID),
		}).
		SetBody(fmt.Sprintf(`{"quantityDelta":%d}`, quantityDelta)).
		Put("/products/{productId}/inventory")

	return responseUpdateCount(response, err)
}
