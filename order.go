package ecwid

import (
	"context"
	"errors"
	"fmt"
)

type (
	// NewOrder https://developers.ecwid.com/api-documentation/orders#create-order
	NewOrder struct {
		Subtotal                        float32            `json:"subtotal"`
		Total                           float32            `json:"total"`
		Email                           string             `json:"email"`
		PaymentMethod                   string             `json:"paymentMethod"`
		PaymentModule                   string             `json:"paymentModule"`
		Tax                             float32            `json:"tax"`
		CustomerTaxExempt               bool               `json:"customerTaxExempt"`
		CustomerTaxID                   string             `json:"customerTaxId"`
		CustomerTaxIDValid              bool               `json:"customerTaxIdValid"`
		ReversedTaxApplied              bool               `json:"reversedTaxApplied"`
		IPAddress                       string             `json:"ipAddress"`
		CouponDiscount                  float32            `json:"couponDiscount"`
		PaymentStatus                   PaymentStatus      `json:"paymentStatus"`
		FulfillmentStatus               FulfillmentStatus  `json:"fulfillmentStatus"`
		RefererURL                      string             `json:"refererUrl"`
		OrderComments                   string             `json:"orderComments"`
		VolumeDiscount                  float32            `json:"volumeDiscount"`
		CustomerID                      uint64             `json:"customerId"`
		Hidden                          bool               `json:"hidden"`
		MembershipBasedDiscount         float32            `json:"membershipBasedDiscount"`
		TotalAndMembershipBasedDiscount float32            `json:"totalAndMembershipBasedDiscount"`
		Discount                        float32            `json:"discount"`
		GlobalReferer                   string             `json:"globalReferer"`
		CreateDate                      DateTime           `json:"createDate"`
		CustomerGroup                   string             `json:"customerGroup"`
		DiscountCoupon                  DiscountCouponInfo `json:"discountCoupon"`
		Items                           []OrderItem        `json:"items"`
		BillingPerson                   PersonInfo         `json:"billingPerson"`
		ShippingPerson                  PersonInfo         `json:"shippingPerson"`
		ShippingOption                  ShippingOptionInfo `json:"shippingOption"`
		HandlingFee                     HandlingFeeInfo    `json:"handlingFee"`
		AdditionalInfo                  map[string]string  `json:"additionalInfo"`
		PaymentParams                   map[string]string  `json:"paymentParams"`
		DiscountInfo                    []DiscountInfo     `json:"discountInfo"`
		TrackingNumber                  string             `json:"trackingNumber"`
		PaymentMessage                  string             `json:"paymentMessage"`
		ExternalTransactionID           string             `json:"externalTransactionId"`
		AffiliateID                     string             `json:"affiliateId"`
		CreditCardStatus                CreditCardStatus   `json:"creditCardStatus"`
		PrivateAdminNotes               string             `json:"privateAdminNotes"`
		PickupTime                      DateTime           `json:"pickupTime"`
		AcceptMarketing                 bool               `json:"acceptMarketing"`
		DisableAllCustomerNotifications bool               `json:"disableAllCustomerNotifications"`
		ExternalFulfillment             bool               `json:"externalFulfillment"`
		ExternalOrderID                 string             `json:"externalOrderId"`
	}

	// Order https://developers.ecwid.com/api-documentation/orders#get-order-details
	Order struct {
		NewOrder
		OrderNumber       uint64           `json:"orderNumber"`
		VendorOrderNumber string           `json:"vendorOrderNumber"`
		USDTotal          float32          `json:"usdTotal"`
		UpdateDate        DateTime         `json:"updateDate"`
		CreateTimestamp   uint64           `json:"createTimestamp"`
		UpdateTimestamp   uint64           `json:"updateTimestamp"`
		CustomerGroupID   uint64           `json:"customerGroupId"`
		PredictedPackages PredictedPackage `json:"predictedPackages"`
		ExtraFields       ExtraFieldsInfo  `json:"extraFields,omitempty"`
		RefundedAmount    float32          `json:"refundedAmount"`
		Refunds           []RefundsInfo    `json:"refunds"`
		RefererID         string           `json:"refererId"`
		TaxesOnShipping   []TaxOnShipping  `json:"taxesOnShipping,omitempty"` // only in Get ???
	}

	// OrdersSearchResponse https://developers.ecwid.com/api-documentation/orders#search-orders
	OrdersSearchResponse struct {
		Total  uint     `json:"total"`
		Count  uint     `json:"count"`
		Offset uint     `json:"offset"`
		Limit  uint     `json:"limit"`
		Orders []*Order `json:"items"`
	}

	// OrderItem contains order items
	OrderItem struct {
		Name                  string            `json:"name"`
		Quantity              uint              `json:"quantity"`
		ProductID             uint64            `json:"productId"`
		CategoryID            uint64            `json:"categoryId"`
		Price                 float32           `json:"price"`
		ProductPrice          float32           `json:"productPrice"`
		Weight                float32           `json:"weight"`
		Sku                   string            `json:"sku"`
		ShortDescription      string            `json:"shortDescription"`
		Tax                   float32           `json:"tax"`
		Shipping              float32           `json:"shipping"`
		QuantityInStock       uint              `json:"quantityInStock"`
		IsShippingRequired    bool              `json:"isShippingRequired"`
		TrackQuantity         bool              `json:"trackQuantity"`
		FixedShippingRateOnly bool              `json:"fixedShippingRateOnly"`
		FixedShippingRate     float32           `json:"fixedShippingRate"`
		Digital               bool              `json:"digital"`
		CouponApplied         bool              `json:"couponApplied"`
		SelectedOptions       []OrderItemOption `json:"selectedOptions"`
		Taxes                 []OrderItemTax    `json:"taxes"`
		Dimensions            ProductDimensions `json:"dimensions"`
	}
)

// OrdersSearch search or filter orders in a store
// filter:
// keywords totalFrom totalTo createdFrom createdTo updatedFrom updatedTo
// couponCode orderuint64 vendorOrderuint64
// email customerId paymentMethod shippingMethod paymentStatus fulfillmentStatus
// acceptMarketing refererId productId offset limit
func (c *Client) OrdersSearch(filter map[string]string) (*OrdersSearchResponse, error) {
	response, err := c.R().
		SetQueryParams(filter).
		Get("/orders")

	var result OrdersSearchResponse
	return &result, responseUnmarshal(response, err, &result)
}

// Orders 'iterable' by filtered store orders
func (c *Client) Orders(ctx context.Context, filter map[string]string) <-chan *Order {
	orderChan := make(chan *Order)

	go func() {
		defer close(orderChan)

		c.OrdersTrampoline(filter, func(index int, order *Order) error {
			// FIXME silent error. maybe orderChan <- nil ?
			select {
			case <-ctx.Done():
				return errors.New("break")
			case orderChan <- order:
			}
			return nil
		})
	}()

	return orderChan
}

// OrdersTrampoline call on each orders
func (c *Client) OrdersTrampoline(filter map[string]string, fn func(int, *Order) error) error {
	filterCopy := make(map[string]string)
	for k, v := range filter {
		filterCopy[k] = v
	}

	for {
		resp, err := c.OrdersSearch(filterCopy)
		if err != nil {
			return err
		}

		for index, order := range resp.Orders {
			if err := fn(index, order); err != nil {
				return err
			}
		}

		if resp.Offset+resp.Count >= resp.Total {
			return nil
		}
		filterCopy["offset"] = fmt.Sprintf("%d", resp.Offset+resp.Count)
	}
}

// OrderGet gets all details of a specific order in an Ecwid store by its ID
func (c *Client) OrderGet(orderID uint64) (*Order, error) {
	response, err := c.R().
		Get(fmt.Sprintf("/orders/%d", orderID))

	var result Order
	return &result, responseUnmarshal(response, err, &result)
}

// TODO add more order methods
// Get order invoice
// Update order
// Delete order
// Create order
// Upload item option file
// Delete item option file
// Delete all item option’s files
