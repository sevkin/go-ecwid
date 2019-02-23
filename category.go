package ecwid

import (
	"fmt"
	"html/template"
)

type (
	// TODO ImageDetails struct {
	// 	URL    string `json:"url"`    //	string	Image URL
	// 	Width  uint   `json:"width"`  //	integer	Image width
	// 	Height uint   `json:"height"` //	integer	Image height
	// }

	// NewCategory https://developers.ecwid.com/api-documentation/categories#add-new-category
	NewCategory struct {
		Name string `json:"name,omitempty"`
		// TODO ParentID can`t set to 0 (if omitempty) or 0 when not set on update
		ParentID    uint64        `json:"parentId,omitempty"` // `json:"parentId"`
		OrderBy     int           `json:"orderBy,omitempty"`
		Description template.HTML `json:"description,omitempty"`
		Enabled     bool          `json:"enabled,omitempty"`
		ProductIDs  []uint64      `json:"productIds,omitempty"`
	}

	// Category https://developers.ecwid.com/api-documentation/categories#get-categories
	Category struct {
		*NewCategory
		ID               uint64 `json:"id"`
		HdThumbnailURL   string `json:"hdThumbnailUrl,omitempty"`
		ThumbnailURL     string `json:"thumbnailUrl,omitempty"`
		ImageURL         string `json:"imageUrl,omitempty"`
		OriginalImageURL string `json:"originalImageUrl,omitempty"`
		// OriginalImage       *ImageDetails `json:"originalImage,omitempty"`
		URL                 string `json:"url,omitempty"`
		ProductCount        uint   `json:"productCount,omitempty"`
		EnabledProductCount uint   `json:"enabledProductCount,omitempty"`
	}

	// GetCategoriesResponse is basic details of found categories
	GetCategoriesResponse struct {
		Total  uint        `json:"total"`
		Count  uint        `json:"count"`
		Offset uint        `json:"offset"`
		Limit  uint        `json:"limit"`
		Items  []*Category `json:"items"`
	}
)

// GetCategories search or filter categories in a store catalog
// The response provides basic details of found categories
func (c *Client) GetCategories(filter map[string]string) (*GetCategoriesResponse, error) {
	// filter:
	// parent number, hidden_categories bool, offset number, limit number,
	// productIds array?, baseUrl string, cleanUrls string

	response, err := c.R().
		SetQueryParams(filter).
		Get("/categories")

	var result GetCategoriesResponse
	return &result, responseUnmarshal(response, err, &result)
}

// GetCategory gets all details of a specific category in an Ecwid store by its ID
func (c *Client) GetCategory(categoryID uint) (*Category, error) {
	response, err := c.R().
		SetPathParams(map[string]string{
			"categoryId": fmt.Sprintf("%d", categoryID),
		}).
		Get("/categories/{categoryId}")

	var result Category
	return &result, responseUnmarshal(response, err, &result)
}

// AddNewCategory creates a new category in an Ecwid store
// returns new categoryId
func (c *Client) AddNewCategory(category *NewCategory) (uint64, error) {
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(category).
		Post("/categories")

	return responseAdd(response, err)
}

// UpdateCategory update an existing category in an Ecwid store referring to its ID
// before update use GetCategory to retrieve full data
func (c *Client) UpdateCategory(categoryID uint64, category *NewCategory) error {
	response, err := c.R().
		SetPathParams(map[string]string{
			"categoryId": fmt.Sprintf("%d", categoryID),
		}).
		SetHeader("Content-Type", "application/json").
		SetBody(category).
		Put("/categories/{categoryId}")

	return responseUpdate(response, err)
}

// DeleteCategory delete a category from an Ecwid store referring to its ID
func (c *Client) DeleteCategory(categoryID uint64) error {
	response, err := c.R().
		SetPathParams(map[string]string{
			"categoryId": fmt.Sprintf("%d", categoryID),
		}).
		Delete("/categories/{categoryId}")

	return responseDelete(response, err)
}
