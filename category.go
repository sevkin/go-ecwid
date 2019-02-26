package ecwid

import (
	"context"
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

	// CategoriesGetResponse is basic details of found categories
	CategoriesGetResponse struct {
		Total  uint        `json:"total"`
		Count  uint        `json:"count"`
		Offset uint        `json:"offset"`
		Limit  uint        `json:"limit"`
		Items  []*Category `json:"items"`
	}
)

// CategoriesGet search or filter categories in a store catalog
// The response provides basic details of found categories
func (c *Client) CategoriesGet(filter map[string]string) (*CategoriesGetResponse, error) {
	// filter:
	// parent number, hidden_categories bool, offset number, limit number,
	// productIds array?, baseUrl string, cleanUrls bool

	response, err := c.R().
		SetQueryParams(filter).
		Get("/categories")

	var result CategoriesGetResponse
	return &result, responseUnmarshal(response, err, &result)
}

// Categories 'iterable' by filtered store categories
func (c *Client) Categories(ctx context.Context, filter map[string]string) <-chan *Category {
	catChan := make(chan *Category)

	filterCopy := make(map[string]string)
	for k, v := range filter {
		filterCopy[k] = v
	}

	go func(filter map[string]string) {
		defer close(catChan)

		for {
			resp, err := c.CategoriesGet(filter)
			if err != nil {
				return
			}

			for _, category := range resp.Items {
				select {
				case <-ctx.Done():
					return
				case catChan <- category:
				}
			}

			if resp.Offset+resp.Count >= resp.Total {
				return
			}
			filter["offset"] = fmt.Sprintf("%d", resp.Offset+resp.Count)
		}
	}(filterCopy)

	return catChan
}

// CategoryGet gets all details of a specific category in an Ecwid store by its ID
func (c *Client) CategoryGet(categoryID uint64) (*Category, error) {
	response, err := c.R().
		Get(fmt.Sprintf("/categories/%d", categoryID))

	var result Category
	return &result, responseUnmarshal(response, err, &result)
}

// CategoryAdd creates a new category in an Ecwid store
// returns new categoryId
func (c *Client) CategoryAdd(category *NewCategory) (uint64, error) {
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(category).
		Post("/categories")

	return responseAdd(response, err)
}

// CategoryUpdate update an existing category in an Ecwid store referring to its ID
// before update use CategoryGet to retrieve full data
func (c *Client) CategoryUpdate(categoryID uint64, category *NewCategory) error {
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(category).
		Put(fmt.Sprintf("/categories/%d", categoryID))

	return responseUpdate(response, err)
}

// CategoryDelete delete a category from an Ecwid store referring to its ID
func (c *Client) CategoryDelete(categoryID uint64) error {
	response, err := c.R().
		Delete(fmt.Sprintf("/categories/%d", categoryID))

	return responseDelete(response, err)
}
