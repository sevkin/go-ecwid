package ecwid

import (
	"encoding/json"
	"errors"
	"fmt"

	resty "gopkg.in/resty.v1"
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
		Total    uint      `json:"total"`
		Count    uint      `json:"count"`
		Offset   uint      `json:"offset"`
		Limit    uint      `json:"limit"`
		Products []Product `json:"items"`
	}
)

func errorResponse(response *resty.Response) error {
	var result struct {
		ErrorMessage string `json:"errorMessage"`
	}
	if err := json.Unmarshal(response.Body(), &result); err == nil && len(result.ErrorMessage) > 0 {
		return errors.New(result.ErrorMessage)
	}
	return errors.New(response.String())
}

// SearchProducts search or filter products in a store catalog
func (c *Client) SearchProducts(filter map[string]string) (*SearchProductsResponse, error) {
	// keyword={keyword}&priceFrom={priceFrom}&priceTo={priceTo}&category={category}&withSubcategories={withSubcategories}&sortBy={sortBy}&offset={offset}&limit={limit}&createdFrom={createdFrom}&createdTo={createdTo}&updatedFrom={updatedFrom}&updatedTo={updatedTo}&enabled={enabled}&inStock={inStock}&field{attributeName}={attributeValues}&field{attributeId}={attributeValues}&sku={sku}&productId={productId}&baseUrl={baseUrl}&cleanUrls={cleanUrls}&onsale={onsale}&option_{optionName}={optionValues}&attribute_{attributeName}={attributeValues}&token={token}

	response, err := c.R().SetQueryParams(filter).Get("/products")
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		return nil, errorResponse(response)
	}

	var result SearchProductsResponse
	err = json.Unmarshal(response.Body(), &result)
	if err == nil {
		return &result, nil
	}
	return nil, err
}

// GetProduct gets all details of a specific product in an Ecwid store by its ID
func (c *Client) GetProduct(productID uint) (*Product, error) {
	response, err := c.R().SetPathParams(map[string]string{
		"productId": fmt.Sprintf("%d", productID),
	}).Get("/products/{productId}")

	if response.StatusCode() != 200 {
		return nil, errorResponse(response)
	}

	var result Product
	err = json.Unmarshal(response.Body(), &result)
	if err == nil {
		return &result, nil
	}
	return nil, err
}

// AddProduct creates a new product in an Ecwid store
// returns new productId
func (c *Client) AddProduct(product *Product) (uint, error) {
	// FIXME!!! похоже на каждый запрос (add|get|search|update) набор полей различный

	body, err := json.Marshal(product)
	if err != nil {
		return 0, err
	}
	response, err := c.R().SetBody(body).Post("/products")

	if response.StatusCode() != 200 { // TODO check real ecwid codes
		return 0, errorResponse(response)
	}

	var result struct {
		ID uint `json:"id"`
	}
	err = json.Unmarshal(response.Body(), &result)
	if err == nil {
		return result.ID, nil
	}
	return 0, err
}

// UpdateProduct update an existing product in an Ecwid store referring to its ID
// before update use GetProduct to retrieve full data
func (c *Client) UpdateProduct(productID uint, product *Product) error {
	body, err := json.Marshal(product)
	if err != nil {
		return err
	}

	response, err := c.R().SetPathParams(map[string]string{
		"productId": fmt.Sprintf("%d", productID),
	}).SetBody(body).Put("/products/{productId}")
	// TODO may be marshal body not need

	if response.StatusCode() != 200 {
		return errorResponse(response)
	}

	var result struct {
		UpdateCount uint `json:"updateCount"`
	}
	err = json.Unmarshal(response.Body(), &result)
	if err == nil {
		if result.UpdateCount == 1 {
			return nil
		}
		return errors.New("no products updated")
	}
	return err

}

// TODO try to pass partial json with help https://github.com/tidwall/sjson
// func (c *Client) UpdateProductJson(productID uint, productJson string) error {

// DeleteProduct delete a product from an Ecwid store referring to its ID
func (c *Client) DeleteProduct(productID uint) error {
	response, err := c.R().SetPathParams(map[string]string{
		"productId": fmt.Sprintf("%d", productID),
	}).Delete("/products/{productId}")

	if response.StatusCode() != 200 {
		return errorResponse(response)
	}

	var result struct {
		UpdateCount uint `json:"deleteCount"`
	}
	err = json.Unmarshal(response.Body(), &result)
	if err == nil {
		if result.UpdateCount == 1 {
			return nil
		}
		return errors.New("no products deleted")
	}
	return err
}

// AdjustProductInventory increase or decrease the product’s stock quantity by a delta quantity
func (c *Client) AdjustProductInventory(productID uint, quantityDelta int) (int, error) {
	response, err := c.R().SetPathParams(map[string]string{
		"productId": fmt.Sprintf("%d", productID),
	}).SetBody(fmt.Sprintf(`{"quantityDelta":%d}`, quantityDelta)).Put("/products/{productId}/inventory")

	if response.StatusCode() != 200 {
		return 0, errorResponse(response)
	}

	var result struct {
		UpdateCount int `json:"updateCount"`
	}
	err = json.Unmarshal(response.Body(), &result)
	if err == nil {
		return result.UpdateCount, nil
	}
	return 0, err
}
