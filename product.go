package ecwid

import (
	"fmt"
	"html/template"
)

type (
	// NewProduct https://developers.ecwid.com/api-documentation/products#add-a-product
	NewProduct struct {
		Name                  string        `json:"name,omitempty"` // mandatory for ProductAdd
		Sku                   string        `json:"sku,omitempty"`
		Quantity              int           `json:"quantity,omitempty"`
		Unlimited             bool          `json:"unlimited,omitempty"`
		Price                 float32       `json:"price,omitempty"`
		CompareToPrice        float32       `json:"compareToPrice,omitempty"`
		IsShippingRequired    bool          `json:"isShippingRequired,omitempty"`
		Weight                float32       `json:"weight,omitempty"`
		ProductClassID        uint64        `json:"productClassId,omitempty"`
		Created               string        `json:"created,omitempty"` // may be in Product not NewProduct
		Enabled               bool          `json:"enabled,omitempty"`
		WarningLimit          uint          `json:"warningLimit,omitempty"`
		FixedShippingRateOnly bool          `json:"fixedShippingRateOnly,omitempty"`
		FixedShippingRate     float32       `json:"fixedShippingRate,omitempty"`
		Description           template.HTML `json:"description,omitempty"`
		SeoTitle              string        `json:"seoTitle,omitempty"`
		SeoDescription        string        `json:"seoDescription,omitempty"`
		DefaultCategoryID     uint64        `json:"defaultCategoryId,omitempty"`
		ShowOnFrontpage       int           `json:"showOnFrontpage,omitempty"`

		CategoryIDs     []uint64           `json:"categoryIds,omitempty"`
		WholesalePrices []WholesalePrice   `json:"wholesalePrices,omitempty"`
		Options         []ProductOption    `json:"options,omitempty"`
		Attributes      []AttributeValue   `json:"attributes,omitempty"`
		Tax             *TaxInfo           `json:"tax,omitempty"`
		Shipping        *ShippingSettings  `json:"shipping,omitempty"`
		RelatedProducts *RelatedProducts   `json:"relatedProducts,omitempty"`
		Dimensions      *ProductDimensions `json:"dimensions,omitempty"`

		Media *ProductMedia `json:"media,omitempty"` // ProductUpdate, ProductGet ProductsSearch
		// GalleryImages []GalleryImage `json:"galleryImages,omitempty"` // only ProductUpdate
	}

	// Product https://developers.ecwid.com/api-documentation/products#get-a-product
	Product struct {
		*NewProduct
		ID                                     uint64  `json:"id"`
		InStock                                bool    `json:"inStock,omitempty"`
		DefaultDisplayedPrice                  float32 `json:"defaultDisplayedPrice,omitempty"`
		DefaultDisplayedPriceFormatted         string  `json:"defaultDisplayedPriceFormatted,omitempty"`
		CompareToPriceFormatted                string  `json:"compareToPriceFormatted,omitempty"`
		CompareToPriceDiscount                 float32 `json:"compareToPriceDiscount,omitempty"`
		CompareToPriceDiscountFormatted        string  `json:"compareToPriceDiscountFormatted,omitempty"`
		CompareToPriceDiscountPercent          float32 `json:"compareToPriceDiscountPercent,omitempty"`
		CompareToPriceDiscountPercentFormatted string  `json:"compareToPriceDiscountPercentFormatted,omitempty"`
		URL                                    string  `json:"url,omitempty"`
		Updated                                string  `json:"updated,omitempty"`
		CreateTimestamp                        uint    `json:"createTimestamp,omitempty"`
		UpdateTimestamp                        uint    `json:"updateTimestamp,omitempty"`
		DefaultCombinationID                   uint64  `json:"defaultCombinationId,omitempty"`
		IsSampleProduct                        bool    `json:"isSampleProduct,omitempty"`

		Combinations []ProductVariation `json:"combinations,omitempty"`
		Categories   []CategoriesInfo   `json:"categories,omitempty"`
		// Files        []ProductFile      `json:"files,omitempty"`
		// Favorites    *FavoritesStats    `json:"favorites,omitempty"`
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
// before update use ProductGet to retrieve full data
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
		SetBody(fmt.Sprintf(`{"quantityDelta":%d}`, quantityDelta)).
		Put(fmt.Sprintf("/products/%d/inventory", productID))

	return responseUpdateCount(response, err)
}
