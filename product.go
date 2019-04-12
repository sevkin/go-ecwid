package ecwid

import (
	"context"
	"errors"
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
		ShowOnFrontpage       int                `json:"showOnFrontpage,omitempty"`
		CategoryIDs           []uint64           `json:"categoryIds,omitempty"`
		WholesalePrices       []WholesalePrice   `json:"wholesalePrices,omitempty"`
		Options               []ProductOption    `json:"options,omitempty"`
		Attributes            Attributes         `json:"attributes,omitempty"`
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
		CreateTimestamp                        uint64             `json:"createTimestamp"`
		UpdateTimestamp                        uint64             `json:"updateTimestamp"`
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

// Products 'iterable' by filtered store products
func (c *Client) Products(ctx context.Context, filter map[string]string) <-chan *Product {
	prodChan := make(chan *Product)

	go func() {
		defer close(prodChan)

		c.ProductsTrampoline(filter, func(index int, product *Product) error {
			// FIXME silent error. maybe prodChan <- nil ?
			select {
			case <-ctx.Done():
				return errors.New("break")
			case prodChan <- product:
			}
			return nil
		})
	}()

	return prodChan
}

// ProductsTrampoline call on each products
func (c *Client) ProductsTrampoline(filter map[string]string, fn func(int, *Product) error) error {
	filterCopy := make(map[string]string)
	for k, v := range filter {
		filterCopy[k] = v
	}

	for {
		resp, err := c.ProductsSearch(filterCopy)
		if err != nil {
			return err
		}

		for index, product := range resp.Products {
			if err := fn(index, product); err != nil {
				return err
			}
		}

		if resp.Offset+resp.Count >= resp.Total {
			return nil
		}
		filterCopy["offset"] = fmt.Sprintf("%d", resp.Offset+resp.Count)
	}
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

	_, err = responseDelete(response, err)
	return err
}

// ProductInventoryAdjust increase or decrease the product’s stock quantity by a delta quantity
func (c *Client) ProductInventoryAdjust(productID uint, quantityDelta int) (int, error) {
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"quantityDelta":%d}`, quantityDelta)).
		Put(fmt.Sprintf("/products/%d/inventory", productID))

	return responseUpdateCount(response, err)
}
