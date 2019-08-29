package ecwid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TypesTestSuite struct {
	suite.Suite
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (suite *TypesTestSuite) TestIDUnmashalJSON() {
	const (
		sidN = `{"id": -1}`
		sidZ = `{"id": 0}`
		sidP = `{"id": 1}`
	)

	var sid struct {
		ID ID
	}

	err := json.Unmarshal([]byte(sidP), &sid)
	suite.Nil(err)
	suite.Equal(ID(1), sid.ID)

	err = json.Unmarshal([]byte(sidZ), &sid)
	suite.Nil(err)
	suite.Equal(ID(0), sid.ID)

	err = json.Unmarshal([]byte(sidN), &sid)
	suite.Nil(err)
	suite.Equal(ID(0), sid.ID)
}

func (suite *TypesTestSuite) TestIDMashalJSON() {
	var sid = struct {
		ID ID `json:"id"`
	}{1}

	data, err := json.Marshal(&sid)
	suite.Nil(err)
	suite.Equal(`{"id":1}`, string(data))
}
