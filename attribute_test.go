package ecwid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttributesUnmarshalJSON(t *testing.T) {
	var av Attributes
	err := json.Unmarshal([]byte(`[{"id":1,"name":"n1","value":"v1"}]`), &av)
	assert.Nil(t, err)
	assert.Equal(t, Attributes{
		attributes: []Attribute{
			Attribute{
				ID:    1,
				Name:  "n1",
				Value: "v1",
			},
		},
	}, av)
}

func TestAttributesMarshalJSON(t *testing.T) {
	av := Attributes{
		attributes: []Attribute{
			Attribute{
				ID:    1,
				Name:  "n1",
				Value: "v1",
			},
		},
	}
	b, err := json.Marshal(&av)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`[{"id":1,"name":"n1","value":"v1"}]`), b)
}

func fixtureAttributes() []Attribute {
	return []Attribute{
		Attribute{
			ID:    1,
			Name:  "n1",
			Value: "v1",
		},
		Attribute{
			ID:    2,
			Name:  "n2",
			Value: "v2",
		},
	}
}

func TestAttributesGetByName(t *testing.T) {
	av := Attributes{
		attributes: fixtureAttributes(),
	}

	// found
	a := av.GetByName("n2")
	assert.NotNil(t, a)
	assert.Equal(t, "v2", a.Value)
	assert.Equal(t, &av.attributes[1], a)
	a.Value = "new2"
	assert.Equal(t, "new2", av.attributes[1].Value)

	// not found
	a = av.GetByName("n3")
	assert.Nil(t, a)
}

func TestAttributesGetByID(t *testing.T) {
	av := Attributes{
		attributes: fixtureAttributes(),
	}

	// found
	a := av.GetByID(2)
	assert.NotNil(t, a)
	assert.Equal(t, "v2", a.Value)
	assert.Equal(t, &av.attributes[1], a)
	a.Value = "new2"
	assert.Equal(t, "new2", av.attributes[1].Value)

	// not found
	a = av.GetByID(3)
	assert.Nil(t, a)
}

func TestAttributesDelete(t *testing.T) {
	av := Attributes{
		attributes: fixtureAttributes(),
	}
	// by ref
	assert.True(t, av.Delete(av.GetByName("n2")))
	assert.Equal(t, 1, len(av.attributes))
	assert.Equal(t, "v1", av.attributes[0].Value)

	// else
	assert.False(t, av.Delete(&Attribute{}))
	assert.Equal(t, 1, len(av.attributes))

	assert.False(t, av.Delete(nil))
	assert.Equal(t, 1, len(av.attributes))
}

func TestAttributesAppend(t *testing.T) {
	av := Attributes{
		attributes: fixtureAttributes(),
	}

	a := av.Append(nil)
	assert.NotNil(t, a)
	assert.Equal(t, 3, len(av.attributes))
	a.Value = "new3"
	assert.Equal(t, "new3", av.attributes[len(av.attributes)-1].Value)

	a = av.Append(&Attribute{
		Name:  "n4",
		Value: "v4",
	})
	assert.NotNil(t, a)
	assert.Equal(t, 4, len(av.attributes))
	a.Value = "new4"
	assert.Equal(t, "new4", av.attributes[len(av.attributes)-1].Value)
}

func TestAttributesCopyTo(t *testing.T) {
	av := Attributes{
		attributes: fixtureAttributes(),
	}
	dest := Attributes{}
	av.CopyTo(&dest)

	assert.Equal(t, av.attributes, dest.attributes)

	dest.attributes[0].Value = "new1"
	assert.NotEqual(t, av.attributes, dest.attributes)
}

func TestAttributesIsEqualTo(t *testing.T) {

	a := Attributes{attributes: fixtureAttributes()}
	a1 := Attributes{attributes: fixtureAttributes()}
	assert.True(t, a.IsEqualTo(&a1))

	a2 := Attributes{attributes: []Attribute{
		Attribute{
			ID:    2,
			Name:  "n2",
			Value: "v2",
		},
		Attribute{
			ID:    1,
			Name:  "n1",
			Value: "v1",
		},
	}}
	assert.True(t, a.IsEqualTo(&a2))

	a3 := Attributes{attributes: fixtureAttributes()}
	a3.attributes[0].ID = 3
	assert.False(t, a.IsEqualTo(&a3))

	a4 := Attributes{attributes: fixtureAttributes()}
	a4.attributes[0].Name = "n3"
	assert.False(t, a.IsEqualTo(&a4))

	a5 := Attributes{attributes: fixtureAttributes()}
	a5.attributes[0].Value = "v3"
	assert.False(t, a.IsEqualTo(&a5))

	a6 := Attributes{attributes: []Attribute{
		Attribute{
			ID:    1,
			Name:  "n1",
			Value: "v1",
		},
	}}
	assert.False(t, a.IsEqualTo(&a6))
}
