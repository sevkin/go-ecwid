package ecwid

import "encoding/json"

type (

	// SearchResponse - common parts of all search funcs
	SearchResponse struct {
		Total  uint `json:"total"`
		Count  uint `json:"count"`
		Offset uint `json:"offset"`
		Limit  uint `json:"limit"`
	}

	// Custom fields

	// ID object identifier.
	// Like a NULL if zero or negative.
	ID uint64

	// DateTime some like "2015-09-20 19:59:43 +0000"
	DateTime string // TODO DateTime => time.Date + Marshal|Unmarshal|String

	// ModifierType - price modifier type
	ModifierType string
)

// ModifierType types
const (
	ModifierPercent  ModifierType = "PERCENT"
	ModifierAbsolute ModifierType = "ABSOLUTE"
)

// UnmarshalJSON unmarshal negative as zero
func (id *ID) UnmarshalJSON(data []byte) error {
	var i int64
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	if i < 0 {
		*id = 0
		return nil
	}
	*id = ID(i)
	return nil
}
