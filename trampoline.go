package ecwid

import "fmt"

// closure hell instead of generics

type searchTrampoliner func(map[string]string, uint) (*SearchResponse, error)

func searchTrampoline(filter map[string]string, trampoliner searchTrampoliner) error {
	filterCopy := make(map[string]string)
	for k, v := range filter {
		filterCopy[k] = v
	}

	index := uint(0)

	for {
		resp, err := trampoliner(filterCopy, index)
		if err != nil {
			return err
		}

		index += resp.Count

		if resp.Offset+resp.Count >= resp.Total {
			return nil
		}
		filterCopy["offset"] = fmt.Sprintf("%d", resp.Offset+resp.Count)
	}
}

// ////////////////////////////////////////////////////////////////////////////

// ProductsTrampoline call on each product
func (c *Client) ProductsTrampoline(filter map[string]string, fn func(uint, *Product) error) error {

	return searchTrampoline(filter, func(filter map[string]string, index uint) (*SearchResponse, error) {
		resp, err := c.ProductsSearch(filter)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			if err := fn(index, item); err != nil {
				return nil, err
			}
			index++
		}

		return &resp.SearchResponse, nil
	})
}

// ////////////////////////////////////////////////////////////////////////////

// CategoriesTrampoline call on each category
func (c *Client) CategoriesTrampoline(filter map[string]string, fn func(uint, *Category) error) error {

	return searchTrampoline(filter, func(filter map[string]string, index uint) (*SearchResponse, error) {
		resp, err := c.CategoriesSearch(filter)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			if err := fn(index, item); err != nil {
				return nil, err
			}
			index++
		}

		return &resp.SearchResponse, nil
	})
}

// ////////////////////////////////////////////////////////////////////////////

// OrdersTrampoline call on each order
func (c *Client) OrdersTrampoline(filter map[string]string, fn func(uint, *Order) error) error {

	return searchTrampoline(filter, func(filter map[string]string, index uint) (*SearchResponse, error) {
		resp, err := c.OrdersSearch(filter)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			if err := fn(index, item); err != nil {
				return nil, err
			}
			index++
		}

		return &resp.SearchResponse, nil
	})
}
