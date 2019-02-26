package ecwid

import (
	"fmt"
	"html/template"
)

type (
	// NewProduct https://developers.ecwid.com/api-documentation/products#add-a-product
	NewProduct struct {
		Name                  string             `json:"name,omitempty"` // mandatory for ProductAdd
		Sku                   string             `json:"sku,omitempty"`
		Quantity              int                `json:"quantity"`
		Unlimited             bool               `json:"unlimited"`
		Price                 float32            `json:"price,omitempty"`
		CompareToPrice        float32            `json:"compareToPrice,omitempty"`
		IsShippingRequired    bool               `json:"isShippingRequired"`
		Weight                float32            `json:"weight,omitempty"`
		ProductClassID        uint64             `json:"productClassId"`
		Created               string             `json:"created,omitempty"`
		Enabled               bool               `json:"enabled"`
		WarningLimit          uint               `json:"warningLimit"`
		FixedShippingRateOnly bool               `json:"fixedShippingRateOnly"`
		FixedShippingRate     float32            `json:"fixedShippingRate,omitempty"`
		Description           template.HTML      `json:"description,omitempty"`
		SeoTitle              string             `json:"seoTitle,omitempty"`
		SeoDescription        string             `json:"seoDescription,omitempty"`
		DefaultCategoryID     uint64             `json:"defaultCategoryId"`
		ShowOnFrontpage       int                `json:"showOnFrontpage"`
		CategoryIDs           []uint64           `json:"categoryIds,omitempty"`
		WholesalePrices       []WholesalePrice   `json:"wholesalePrices,omitempty"`
		Options               []ProductOption    `json:"options,omitempty"`
		Attributes            []AttributeValue   `json:"attributes,omitempty"`
		Tax                   *TaxInfo           `json:"tax"`
		Shipping              *ShippingSettings  `json:"shipping"`
		RelatedProducts       *RelatedProducts   `json:"relatedProducts"`
		Dimensions            *ProductDimensions `json:"dimensions"`
		Media                 *ProductMedia      `json:"media"` // ProductUpdate, ProductGet ProductsSearch
		// GalleryImages         []GalleryImage     `json:"galleryImages,omitempty"` // only ProductUpdate
	}

	// Product https://developers.ecwid.com/api-documentation/products#get-a-product
	Product struct {
		NewProduct
		ID                                     uint64             `json:"id"`
		InStock                                bool               `json:"inStock"`
		DefaultDisplayedPrice                  float32            `json:"defaultDisplayedPrice"`
		DefaultDisplayedPriceFormatted         string             `json:"defaultDisplayedPriceFormatted"`
		CompareToPriceFormatted                string             `json:"compareToPriceFormatted"`
		CompareToPriceDiscount                 float32            `json:"compareToPriceDiscount"`
		CompareToPriceDiscountFormatted        string             `json:"compareToPriceDiscountFormatted"`
		CompareToPriceDiscountPercent          float32            `json:"compareToPriceDiscountPercent"`
		CompareToPriceDiscountPercentFormatted string             `json:"compareToPriceDiscountPercentFormatted"`
		URL                                    string             `json:"url"`
		Updated                                string             `json:"updated"`
		CreateTimestamp                        uint               `json:"createTimestamp"`
		UpdateTimestamp                        uint               `json:"updateTimestamp"`
		DefaultCombinationID                   uint64             `json:"defaultCombinationId"`
		IsSampleProduct                        bool               `json:"isSampleProduct"`
		Combinations                           []ProductVariation `json:"combinations"`
		Categories                             []CategoriesInfo   `json:"categories"`
		// Files                                  []ProductFile      `json:"files"`
		// Favorites                              *FavoritesStats    `json:"favorites"`
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
		Get(fmt.Sprintf("/products/%d", productID))

	var result Product
	return &result, responseUnmarshal(response, err, &result)
}

// ProductAdd creates a new product in an Ecwid store
// returns new productId
func (c *Client) ProductAdd(product *NewProduct) (uint64, error) {
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(product).
		Post("/products")

	return responseAdd(response, err)
}

// ProductUpdate update an existing product in an Ecwid store referring to its ID
func (c *Client) ProductUpdate(productID uint64, product *NewProduct) error {
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(product).
		Put(fmt.Sprintf("/products/%d", productID))

	return responseUpdate(response, err)
}

// TODO try to pass partial json with help https://github.com/tidwall/sjson
// func (c *Client) UpdateProductJson(productID uint, productJson string) error {

// ProductDelete delete a product from an Ecwid store referring to its ID
func (c *Client) ProductDelete(productID uint64) error {
	response, err := c.R().
		Delete(fmt.Sprintf("/products/%d", productID))

	return responseDelete(response, err)
}

// ProductInventoryAdjust increase or decrease the product’s stock quantity by a delta quantity
func (c *Client) ProductInventoryAdjust(productID uint, quantityDelta int) (int, error) {
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"quantityDelta":%d}`, quantityDelta)).
		Put(fmt.Sprintf("/products/%d/inventory", productID))

	return responseUpdateCount(response, err)
}
