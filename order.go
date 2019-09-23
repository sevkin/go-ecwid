package ecwid

import (
	"context"
	"errors"
	"fmt"
)

type (
	// NewOrder https://developers.ecwid.com/api-documentation/orders#create-order
	NewOrder struct {
		Subtotal                        float32             `json:"subtotal,omitempty"`
		Total                           float32             `json:"total,omitempty"`
		Email                           string              `json:"email,omitempty"`
		PaymentMethod                   string              `json:"paymentMethod,omitempty"`
		PaymentModule                   string              `json:"paymentModule,omitempty"`
		Tax                             float32             `json:"tax,omitempty"`
		CustomerTaxExempt               bool                `json:"customerTaxExempt,omitempty"`
		CustomerTaxID                   string              `json:"customerTaxId,omitempty"`
		CustomerTaxIDValid              bool                `json:"customerTaxIdValid,omitempty"`
		ReversedTaxApplied              bool                `json:"reversedTaxApplied,omitempty"`
		IPAddress                       string              `json:"ipAddress,omitempty"`
		CouponDiscount                  float32             `json:"couponDiscount,omitempty"`
		PaymentStatus                   PaymentStatus       `json:"paymentStatus,omitempty"`
		FulfillmentStatus               FulfillmentStatus   `json:"fulfillmentStatus,omitempty"`
		RefererURL                      string              `json:"refererUrl,omitempty"`
		OrderComments                   string              `json:"orderComments,omitempty"`
		VolumeDiscount                  float32             `json:"volumeDiscount,omitempty"`
		CustomerID                      ID                  `json:"customerId,omitempty"`
		Hidden                          bool                `json:"hidden,omitempty"`
		MembershipBasedDiscount         float32             `json:"membershipBasedDiscount,omitempty"`
		TotalAndMembershipBasedDiscount float32             `json:"totalAndMembershipBasedDiscount,omitempty"`
		Discount                        float32             `json:"discount,omitempty"`
		GlobalReferer                   string              `json:"globalReferer,omitempty"`
		CreateDate                      DateTime            `json:"createDate,omitempty"`
		CustomerGroup                   string              `json:"customerGroup,omitempty"`
		DiscountCoupon                  *DiscountCouponInfo `json:"discountCoupon,omitempty"`
		Items                           []*OrderItem        `json:"items,omitempty"`
		BillingPerson                   *PersonInfo         `json:"billingPerson,omitempty"`
		ShippingPerson                  *PersonInfo         `json:"shippingPerson,omitempty"`
		ShippingOption                  *ShippingOptionInfo `json:"shippingOption,omitempty"`
		HandlingFee                     *HandlingFeeInfo    `json:"handlingFee,omitempty"`
		AdditionalInfo                  map[string]string   `json:"additionalInfo,omitempty"`
		PaymentParams                   map[string]string   `json:"paymentParams,omitempty"`
		DiscountInfo                    []*DiscountInfo     `json:"discountInfo,omitempty"`
		TrackingNumber                  string              `json:"trackingNumber,omitempty"`
		PaymentMessage                  string              `json:"paymentMessage,omitempty"`
		ExternalTransactionID           string              `json:"externalTransactionId,omitempty"`
		AffiliateID                     string              `json:"affiliateId,omitempty"`
		CreditCardStatus                *CreditCardStatus   `json:"creditCardStatus,omitempty"`
		PrivateAdminNotes               string              `json:"privateAdminNotes,omitempty"`
		PickupTime                      DateTime            `json:"pickupTime,omitempty"`
		AcceptMarketing                 bool                `json:"acceptMarketing,omitempty"`
		DisableAllCustomerNotifications bool                `json:"disableAllCustomerNotifications,omitempty"`
		ExternalFulfillment             bool                `json:"externalFulfillment,omitempty"`
		ExternalOrderID                 string              `json:"externalOrderId,omitempty"`
	}

	// Order https://developers.ecwid.com/api-documentation/orders#get-order-details
	Order struct {
		NewOrder
		OrderID           ID               `json:"orderNumber"`
		VendorOrderNumber string           `json:"vendorOrderNumber"`
		USDTotal          float32          `json:"usdTotal"`
		UpdateDate        DateTime         `json:"updateDate"`
		CreateTimestamp   uint64           `json:"createTimestamp"`
		UpdateTimestamp   uint64           `json:"updateTimestamp"`
		CustomerGroupID   ID               `json:"customerGroupId"`
		PredictedPackages PredictedPackage `json:"predictedPackages"`
		ExtraFields       ExtraFieldsInfo  `json:"extraFields,omitempty"`
		RefundedAmount    float32          `json:"refundedAmount"`
		Refunds           []RefundsInfo    `json:"refunds"`
		RefererID         string           `json:"refererId"`
		TaxesOnShipping   []TaxOnShipping  `json:"taxesOnShipping,omitempty"` // only in Get ???
	}

	// OrdersSearchResponse https://developers.ecwid.com/api-documentation/orders#search-orders
	OrdersSearchResponse struct {
		SearchResponse
		Items []*Order `json:"items"`
	}

	// OrderItem contains order items
	OrderItem struct {
		Name                  string            `json:"name"`
		Quantity              uint              `json:"quantity"`
		ProductID             ID                `json:"productId"`
		CategoryID            ID                `json:"categoryId"`
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
// couponCode orderId vendorOrderId
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

		c.OrdersTrampoline(filter, func(index uint, order *Order) error {
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

// OrderGet gets all details of a specific order in an Ecwid store by its ID
func (c *Client) OrderGet(orderID ID) (*Order, error) {
	response, err := c.R().
		Get(fmt.Sprintf("/orders/%d", orderID))

	var result Order
	return &result, responseUnmarshal(response, err, &result)
}

// OrderUpdate update an existing order in an Ecwid store referring to its ID
func (c *Client) OrderUpdate(orderID ID, order *NewOrder) error {
	response, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(order).
		Put(fmt.Sprintf("/orders/%d", orderID))

	return responseUpdate(response, err)
}

// TODO add more order methods
// Get order invoice
// Update order
// Delete order
// Create order
// Upload item option file
// Delete item option file
// Delete all item option’s files
