package ecwid

import "fmt"

type (
	// NewProductVariation ...
	NewProductVariation struct {
		Sku                string  `json:"sku,omitempty"`
		Quantity           uint    `json:"quantity"`
		Unlimited          bool    `json:"unlimited"`
		Price              float32 `json:"price,omitempty"`
		Weight             float32 `json:"weight,omitempty"`
		WarningLimit       uint    `json:"warningLimit"`
		CompareToPrice     float32 `json:"compareToPrice,omitempty"`
		IsShippingRequired bool    `json:"isShippingRequired"`

		Options         []OptionValue    `json:"options,omitempty"`
		WholesalePrices []WholesalePrice `json:"wholesalePrices,omitempty"`
		Attributes      Attributes       `json:"attributes,omitempty"`
	}

	// ProductVariation ...
	ProductVariation struct {
		*NewProductVariation
		ID                ID   `json:"id"`
		CombinationNumber uint `json:"combinationNumber"`

		ThumbnailURL      string `json:"thumbnailUrl,omitempty"`
		ImageURL          string `json:"imageUrl,omitempty"`
		SmallThumbnailURL string `json:"smallThumbnailUrl,omitempty"`
		HdThumbnailURL    string `json:"hdThumbnailUrl,omitempty"`
		OriginalImageURL  string `json:"originalImageUrl,omitempty"`
	}
)

// ProductVariationsGet all variations of a specific product in an Ecwid store by its ID
func (c *Client) ProductVariationsGet(productID ID) ([]ProductVariation, error) {
	response, err := c.R().
		Get(fmt.Sprintf("/products/%d/combinations", productID))

	var result []ProductVariation
	return result, responseUnmarshal(response, err, &result)
}

// ProductVariationGet a specific product variation details referring to its ID
func (c *Client) ProductVariationGet(productID, variationID ID) (*ProductVariation, error) {
	response, err := c.R().
		Get(fmt.Sprintf("/products/%d/combinations/%d", productID, variationID))

	var result ProductVariation
	return &result, responseUnmarshal(response, err, &result)
}

// ProductVariationUpdate update a specific product variation details referring to its ID
func (c *Client) ProductVariationUpdate(productID, variationID ID, productVariation *NewProductVariation) error {
	// TODO implement checkLowStockNotification
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(productVariation).
		Put(fmt.Sprintf("/products/%d/combinations/%d", productID, variationID))

	return responseUpdate(response, err)
}

// TODO implement product variations API
// func (c *Client) ProductVariationAdd
// func (c *Client) ProductVariationDelete
// func (c *Client) ProductVariationsDelete
// func (c *Client) ProductVariationInventoryAdjust

// TODO implement product variations image API
// func (c *Client) ProductVariationImageUpload
// func (c *Client) ProductVariationImageUploadFile
// func (c *Client) ProductVariationImageDelete
