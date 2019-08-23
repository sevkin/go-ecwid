package ecwid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AttributeTestSuite struct {
	suite.Suite
}

func TestAttributeTestSuite(t *testing.T) {
	suite.Run(t, new(AttributeTestSuite))
}

func (suite *AttributeTestSuite) TestAttributesUnmarshalJSON() {
	var av Attributes
	err := json.Unmarshal([]byte(`[{"id":1,"name":"n1","value":"v1"}]`), &av)
	suite.Nil(err)
	suite.Equal(Attributes{
		attributes: []Attribute{
			Attribute{
				ID:    1,
				Name:  "n1",
				Value: "v1",
			},
		},
	}, av)
}

func (suite *AttributeTestSuite) TestAttributesMarshalJSON() {
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
	suite.Nil(err)
	suite.Equal([]byte(`[{"id":1,"name":"n1","value":"v1"}]`), b)
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

func (suite *AttributeTestSuite) TestAttributesGetByName() {
	av := Attributes{
		attributes: fixtureAttributes(),
	}

	// found
	a := av.GetByName("n2")
	suite.NotNil(a)
	suite.Equal("v2", a.Value)
	suite.Equal(&av.attributes[1], a)
	a.Value = "new2"
	suite.Equal("new2", av.attributes[1].Value)

	// not found
	a = av.GetByName("n3")
	suite.Nil(a)
}

func (suite *AttributeTestSuite) TestAttributesGetByID() {
	av := Attributes{
		attributes: fixtureAttributes(),
	}

	// found
	a := av.GetByID(2)
	suite.NotNil(a)
	suite.Equal("v2", a.Value)
	suite.Equal(&av.attributes[1], a)
	a.Value = "new2"
	suite.Equal("new2", av.attributes[1].Value)

	// not found
	a = av.GetByID(3)
	suite.Nil(a)
}

func (suite *AttributeTestSuite) TestAttributesDelete() {
	av := Attributes{
		attributes: fixtureAttributes(),
	}
	// by ref
	suite.True(av.Delete(av.GetByName("n2")))
	suite.Equal(1, len(av.attributes))
	suite.Equal("v1", av.attributes[0].Value)

	// else
	suite.False(av.Delete(&Attribute{}))
	suite.Equal(1, len(av.attributes))

	suite.False(av.Delete(nil))
	suite.Equal(1, len(av.attributes))
}

func (suite *AttributeTestSuite) TestAttributesAppend() {
	av := Attributes{
		attributes: fixtureAttributes(),
	}

	a := av.Append(nil)
	suite.NotNil(a)
	suite.Equal(3, len(av.attributes))
	a.Value = "new3"
	suite.Equal("new3", av.attributes[len(av.attributes)-1].Value)

	a = av.Append(&Attribute{
		Name:  "n4",
		Value: "v4",
	})
	suite.NotNil(a)
	suite.Equal(4, len(av.attributes))
	a.Value = "new4"
	suite.Equal("new4", av.attributes[len(av.attributes)-1].Value)
}

func (suite *AttributeTestSuite) TestAttributesCopyTo() {
	av := Attributes{
		attributes: fixtureAttributes(),
	}
	dest := Attributes{}
	av.CopyTo(&dest)

	suite.Equal(av.attributes, dest.attributes)

	dest.attributes[0].Value = "new1"
	suite.NotEqual(av.attributes, dest.attributes)
}

func (suite *AttributeTestSuite) TestAttributesIsEqualTo() {

	a := Attributes{attributes: fixtureAttributes()}
	a1 := Attributes{attributes: fixtureAttributes()}
	suite.True(a.IsEqualTo(&a1))

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
	suite.True(a.IsEqualTo(&a2))

	a3 := Attributes{attributes: fixtureAttributes()}
	a3.attributes[0].ID = 3
	suite.False(a.IsEqualTo(&a3))

	a4 := Attributes{attributes: fixtureAttributes()}
	a4.attributes[0].Name = "n3"
	suite.False(a.IsEqualTo(&a4))

	a5 := Attributes{attributes: fixtureAttributes()}
	a5.attributes[0].Value = "v3"
	suite.False(a.IsEqualTo(&a5))

	a6 := Attributes{attributes: []Attribute{
		Attribute{
			ID:    1,
			Name:  "n1",
			Value: "v1",
		},
	}}
	suite.False(a.IsEqualTo(&a6))
}
