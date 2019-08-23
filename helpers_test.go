package ecwid

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

const ( // consts used in all _test.go files
	storeID = 666
	token   = "token"
)

type ClientTestSuite struct {
	suite.Suite
	client *Client
}

func (suite *ClientTestSuite) SetupTest() {
	suite.client = New(storeID, token)
	httpmock.ActivateNonDefault(suite.client.GetClient())
}

func (suite *ClientTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}
