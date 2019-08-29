package ecwid

import (
	"encoding/json"
	"reflect"
)

type (
	// too bad design for Attribute //

	// for ProductAdd, ProductUpdate
	// ID    ID
	// Alias string
	// Value string

	// for ProductGet, ProductSearch
	// ID    ID
	// Name  string
	// Value string
	// Type  string
	// Show  string

	// for ProductTypeGet, ProductTypesGet
	// ID    ID
	// Name  string
	// Type  string
	// Show  string

	// for ProductTypeAdd
	// Name  string
	// Type  string
	// Show  string

	// Attribute (or AttributeValue in Ecwid) like a key -> value
	Attribute struct {
		ID    ID     `json:"id,omitempty"`    // ID cannot be set to 0 if 'omitempty', but ID is always not 0
		Alias string `json:"alias,omitempty"` // Alias for system attributes like UPC or Brand
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
		Type  string `json:"type,omitempty"` // Type is one of CUSTOM, UPC, BRAND, GENDER, AGE_GROUP, COLOR, SIZE, PRICE_PER_UNIT, UNITS_IN_PRODUCT
		Show  string `json:"show,omitempty"` // Show is ine of NOTSHOW, DESCR, PRICE
	}

	// Attributes just a wrapper for []Attribute with Get, Delete, Append, Copy, Compare
	Attributes struct {
		attributes []Attribute // TODO marshal as array but can hold as map
	}
)

// UnmarshalJSON unmarshal as []Attribute
func (av *Attributes) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &av.attributes)
}

// MarshalJSON marshal as []Attribute
func (av *Attributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(&av.attributes)
}

// GetByName attribute pointer by name.
// returns nil if not found or *Attribute else
func (av *Attributes) GetByName(name string) *Attribute {
	for i, r := range av.attributes {
		if r.Name == name { // FIXME if same name return 1st
			return &av.attributes[i]
		}
	}

	return nil
}

// GetByID is same as Get, just by ID not by Name
func (av *Attributes) GetByID(id ID) *Attribute {
	for i, r := range av.attributes {
		if r.ID == id { // FIXME if same id return 1st
			return &av.attributes[i]
		}
	}

	return nil
}

// Delete deletes Attribute from Attributes
// return true if found and deleted
func (av *Attributes) Delete(a *Attribute) bool {
	if a == nil {
		return false
	}

	for i := range av.attributes {
		if &av.attributes[i] == a {
			av.attributes = append(av.attributes[:i], av.attributes[i+1:]...)
			return true
		}
	}

	return false
}

// Append Attribute to Attributes and returns its pointer
func (av *Attributes) Append(a *Attribute) *Attribute {
	if a == nil {
		av.attributes = append(av.attributes, Attribute{})
	} else {
		av.attributes = append(av.attributes, *a)
	}

	return &av.attributes[len(av.attributes)-1]
}

// CopyTo just make a copy of []Attribute
func (av *Attributes) CopyTo(dest *Attributes) {
	dest.attributes = make([]Attribute, len(av.attributes))
	copy(dest.attributes, av.attributes)
}

// IsEqualTo compares all []Attribute values without its order
func (av *Attributes) IsEqualTo(dest *Attributes) bool {
	aLen := len(av.attributes)
	if aLen != len(dest.attributes) {
		return false
	}

	visited := make([]bool, aLen)
	for i := 0; i < aLen; i++ {
		a := &av.attributes[i]
		found := false
		for j := 0; j < aLen; j++ {
			if visited[j] {
				continue
			}
			if reflect.DeepEqual(*a, dest.attributes[j]) {
				visited[j] = true
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
