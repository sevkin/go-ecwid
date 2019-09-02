package ecwid

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type OrderTestSuite struct {
	ClientTestSuite
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (suite *OrderTestSuite) TestOrdersSearchRequest() {
	expectedEndpoint := fmt.Sprintf(endpoint+"/orders", storeID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			values := req.URL.Query()
			suite.Equal("test order", values.Get("keyword"), "keyword")
			suite.Equal("5", values.Get("limit"), "limit")

			return httpmock.NewStringResponse(200, ""), nil
		})

	suite.client.OrdersSearch(map[string]string{
		"keyword": "test order",
		"limit":   "5",
	})
	suite.Truef(requested, "request failed")
}

func (suite *OrderTestSuite) TestOrders() {
	requestCount := 0

	expected := []string{"o@n.e", "t@w.o", "t@r.ee"}

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			if requestCount < len(expected) {

				requestCount++

				values := req.URL.Query()
				offset, _ := strconv.ParseUint(values.Get("offset"), 10, 64)
				limit, _ := strconv.ParseUint(values.Get("limit"), 10, 64)

				suite.Equal(uint64(requestCount), offset)
				suite.Equal(uint64(1), limit)

				return httpmock.NewJsonResponse(200, OrdersSearchResponse{
					SearchResponse: SearchResponse{
						Total:  3,
						Count:  1,
						Offset: uint(offset),
						Limit:  uint(limit),
					},
					Orders: []*Order{
						&Order{
							NewOrder: NewOrder{
								Email: expected[offset],
							},
						},
					},
				})
			}
			return httpmock.NewStringResponse(400, "too many"), nil
		})

	actual := make([]string, 0, 3)

	filter := map[string]string{
		"offset": "1",
		"limit":  "1",
	}

	for order := range suite.client.Orders(context.Background(),
		filter) {
		actual = append(actual, order.Email)
	}
	suite.Equal(len(expected)-1, requestCount)
	suite.Equal(expected[1:], actual)

	suite.Equal("1", filter["offset"], "filter map must be unchanged")
}

func (suite *OrderTestSuite) TestOrderGet() {
	const (
		orderID ID = 999
	)

	expectedEndpoint := fmt.Sprintf(endpoint+"/orders/%d", storeID, orderID)
	requested := false

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {
			requested = true

			suite.Equal("GET", req.Method, "request method")
			actualEndpoint := strings.Split(req.URL.String(), "?")[0]
			suite.Equal(expectedEndpoint, actualEndpoint, "endpoint")

			return httpmock.NewStringResponse(200,
				fmt.Sprintf(`{"orderNumber":%d, "email":"s@k.y", 
				"items": [{"name": "test"}]}`, orderID)), nil
		})

	o, err := suite.client.OrderGet(orderID)
	suite.Truef(requested, "request failed")

	suite.Nil(err)
	suite.Equal(orderID, o.OrderID, "id")
	suite.Equal("s@k.y", o.Email, "email")
	suite.Equal(1, len(o.Items), "items")
	suite.Equal("test", o.Items[0].Name, "item name")
}
