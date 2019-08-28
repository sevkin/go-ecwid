package ecwid

type (
	// NewOrder fields

	// DiscountCouponInfo contains information about applied coupon
	DiscountCouponInfo struct {
		Name             string                     `json:"name"`             // Coupon title in store control panel
		Code             string                     `json:"code"`             // Coupon code
		Type             DiscountCouponType         `json:"discountType"`     // Discount type: ABS, PERCENT, SHIPPING, ABS_AND_SHIPPING, PERCENT_AND_SHIPPING
		Status           DiscountCouponStatus       `json:"status"`           // Discount coupon state: ACTIVE, PAUSED, EXPIRED or USEDUP
		Discount         float32                    `json:"discount"`         // Discount amount
		LaunchDate       string                     `json:"launchDate"`       // The date of coupon launch, e.g. 2014-06-06 08:00:00 +0000
		ExpirationDate   string                     `json:"expirationDate"`   // Coupon expiration date, e.g. 2014-06-06 08:00:00 +0000
		TotalLimit       float32                    `json:"totalLimit"`       // The minimum order subtotal the coupon applies to
		UsesLimit        DiscountCouponUseLimit     `json:"usesLimit"`        // Number of uses limitation: UNLIMITED, ONCEPERCUSTOMER, SINGLE
		ApplicationLimit string                     `json:"applicationLimit"` // Application limit for discount coupons. Possible values: "UNLIMITED", "NEW_CUSTOMER_ONLY", "REPEAT_CUSTOMER_ONLY"
		CreationDate     string                     `json:"creationDate"`     // Coupon creation date
		OrderCount       uint64                     `json:"orderCount"`       // Number of uses
		CatalogLimit     DiscountCouponCatalogLimit `json:"catalogLimit"`     // Products and categories the coupon can be applied to
	}

	// DiscountCouponCatalogLimit contains products and categories IDs the coupon can be applied to
	DiscountCouponCatalogLimit struct {
		ProductIDs  []uint64 `json:"products"`
		CategoryIDs []uint64 `json:"categories"`
	}

	// PersonInfo contains name and address of the customer
	PersonInfo struct {
		Name                string `json:"name"`                // Full name
		CompanyName         string `json:"companyName"`         // Company name
		Street              string `json:"street"`              // Address line 1 and address line 2, separated by ’\n’
		City                string `json:"city"`                // City
		CountryCode         string `json:"countryCode"`         // Two-letter country code
		CountryName         string `json:"countryName"`         // Country name
		PostalCode          string `json:"postalCode"`          // Postal/ZIP code
		StateOrProvinceCode string `json:"stateOrProvinceCode"` // State code, e.g. NY
		StateOrProvinceName string `json:"stateOrProvinceName"` // State/province name, e.g. New York
		Phone               string `json:"phone"`               // Phone number
	}

	// ShippingOptionInfo contains information about selected shipping option
	ShippingOptionInfo struct {
		ShippingCarrierName  string  `json:"shippingCarrierName"`  // Optional. Is present for orders made with carriers, e.g. USPS or shipping applications.
		ShippingMethodName   string  `json:"shippingMethodName"`   // Shipping option name
		ShippingRate         float32 `json:"shippingRate"`         // Rate
		EstimatedTransitTime string  `json:"estimatedTransitTime"` // Delivery time estimation. Possible formats: number “5”, several days estimate “4-9”
		IsPickup             bool    `json:"isPickup"`             // true if selected shipping option is local pickup. false otherwise
		PickupInstruction    string  `json:"pickupInstruction"`    // Instruction for customer on how to receive their products
	}

	// HandlingFeeInfo contains handling fee details
	HandlingFeeInfo struct {
		Name        string  `json:"name"`        //	Handling fee name set by store admin. E.g. Wrapping
		Value       float32 `json:"value"`       //	Handling fee value
		Description string  `json:"description"` //	Handling fee description for customer
	}

	// DiscountInfo contains information about applied discounts (coupons are not included)
	DiscountInfo struct {
		Value       float32          `json:"value"`       // Discount value
		Type        DiscountInfoType `json:"type"`        // Discount type: ABS or PERCENT
		Base        DiscountInfoBase `json:"base"`        // Discount base, one of ON_TOTAL, ON_MEMBERSHIP, ON_TOTAL_AND_MEMBERSHIP, CUSTOM
		OrderTotal  float32          `json:"orderTotal"`  // Minimum order subtotal the discount applies to
		Description string           `json:"description"` // Description of a discount (for discounts with base == CUSTOM)
	}

	// CreditCardStatus contains status of credit card payment
	CreditCardStatus struct {
		AVSMessage string `json:"avsMessage"` // Address verification status returned by the payment system.
		CVVMessage string `json:"cvvMessage"` // Credit card verification status returned by the payment system.
	}

	// Order fields

	// PredictedPackage contains predicted information about the package to ship items in to customer
	PredictedPackage struct {
		ProductDimensions
		Weight        float32 `json:"weight"`        // Total weight of a predicted package
		DeclaredValue float32 `json:"declaredValue"` // Declared value of a predicted package (subtotal of items in package)
	}

	// ExtraFieldsInfo Additional optional information about order.
	// Total storage of extra fields cannot exceed 8Kb
	ExtraFieldsInfo map[string]string

	// RefundsInfo contains description of all refunds made to order
	RefundsInfo struct {
		Date   DateTime `json:"date"`   //The date/time of a refund, e.g 2014-06-06 18:57:19 +0000
		Source string   `json:"source"` //What action triggered refund. Possible values: "CP" - changed my merchant in Ecwid CP, "API" - changed by another app, "External" - refund made from payment processor website
		Reason string   `json:"reason"` //A text reason for a refund. 256 characters max
		Amount float32  `json:"amount"` //Amount of this specific refund (not total amount refunded for order. see redundedAmount field)
	}

	// TaxOnShipping taxes applied to shipping 'as is’. null for old orders,
	// [] for orders with taxes applied to subtotal only.
	// Are not recalculated if order is updated later manually.
	// Is calculated like: (shippingRate + handlingFee)*(taxValue/100)
	TaxOnShipping struct {
		Name  string  `json:"name"`
		Value float32 `json:"value"`
		Total float32 `json:"total"`
	}

	// OrderItem fields

	// OrderItemOption contains product options values selected by the customer
	OrderItemOption struct {
		Name       string                `json:"name"`
		Type       string                `json:"type"`
		Value      string                `json:"value"`
		Files      []OrderItemOptionFile `json:"files"`
		Selections []SelectionInfo       `json:"selections"`
	}

	// OrderItemOptionFile contains Attached files if OrderOptionType is FILES
	OrderItemOptionFile struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
		Size uint64 `json:"size"`
		URL  string `json:"url"`
	}

	// SelectionInfo contains details of selected product options.
	// If sent in update order request, other fields will be regenerated based on information in this field
	SelectionInfo struct {
		SelectionTitle        string                `json:"selectionTitle"`
		SelectionModifier     float32               `json:"selectionModifier"` // Money or Percent
		SelectionModifierType SelectionModifierType `json:"selectionModifierType"`
	}

	// OrderItemTax - taxes applied to this order item
	OrderItemTax struct {
		Name                    string  `json:"name"`
		Value                   float32 `json:"value"`
		Total                   float32 `json:"total"`
		TaxOnDiscountedSubtotal float32 `json:"taxOnDiscountedSubtotal"`
		TaxOnShipping           float32 `json:"taxOnShipping"`
		IncludeInPrice          bool    `json:"includeInPrice"`
	}
)

// Custom fields

type (
	// DiscountCouponType ABS, PERCENT, SHIPPING, ABS_AND_SHIPPING, PERCENT_AND_SHIPPING
	DiscountCouponType string

	//DiscountCouponStatus ACTIVE, PAUSED, EXPIRED or USEDUP
	DiscountCouponStatus string

	// DiscountCouponUseLimit UNLIMITED, ONCEPERCUSTOMER, SINGLE
	DiscountCouponUseLimit string

	// DiscountInfoBase ABS or PERCENT
	DiscountInfoBase string

	// DiscountInfoType ON_TOTAL, ON_MEMBERSHIP, ON_TOTAL_AND_MEMBERSHIP, CUSTOM
	DiscountInfoType string

	// PaymentStatus type alias to string
	PaymentStatus string

	// FulfillmentStatus type alias to string
	FulfillmentStatus string

	// OrderOptionType of OrderItemOption
	OrderOptionType string

	// SelectionModifierType - price modifier type
	SelectionModifierType string
)

// DiscountCouponType types
const (
	DiscountCouponAbs                DiscountCouponType = "ABS"
	DiscountCouponPercent            DiscountCouponType = "PERCENT"
	DiscountCouponShipping           DiscountCouponType = "SHIPPING"
	DiscountCouponAbsAndShipping     DiscountCouponType = "ABS_AND_SHIPPING"
	DiscountCouponPercentAndShipping DiscountCouponType = "PERCENT_AND_SHIPPING"
)

// DiscountCouponStatus statuses
const (
	DiscountCouponActive  DiscountCouponStatus = "ACTIVE"
	DiscountCouponPaused  DiscountCouponStatus = "PAUSED"
	DiscountCouponExpired DiscountCouponStatus = "EXPIRED"
	DiscountCouponUsedup  DiscountCouponStatus = "USEDUP"
)

// DiscountCouponUseLimit limits
const (
	DiscountCouponUnlimited       DiscountCouponUseLimit = "UNLIMITED"
	DiscountCouponOncePerCustomer DiscountCouponUseLimit = "ONCEPERCUSTOMER"
	DiscountCouponSingle          DiscountCouponUseLimit = "SINGLE"
)

// DiscountInfoBase bases
const (
	DiscountInfoAbs     DiscountInfoBase = "ABS"
	DiscountInfoPercent DiscountInfoBase = "PERCENT"
)

// DiscountInfoType types
const (
	DiscountInfoOnTotal              DiscountInfoType = "ON_TOTAL"
	DiscountInfoOnMembership         DiscountInfoType = "ON_MEMBERSHIP"
	DiscountInfoOnTotalAndMembership DiscountInfoType = "ON_TOTAL_AND_MEMBERSHIP"
	DiscountInfoCustom               DiscountInfoType = "CUSTOM"
)

// Payment statuses
const (
	PaymentAwaiting          PaymentStatus = "AWAITING_PAYMENT"
	PaymentPaid              PaymentStatus = "PAID"
	PaymentCancelled         PaymentStatus = "CANCELLED"
	PaymentRefunded          PaymentStatus = "REFUNDED"
	PaymentPartiallyRefunded PaymentStatus = "PARTIALLY_REFUNDED"
	PaymentIncomplete        PaymentStatus = "INCOMPLETE"
)

// Fulfillment statuses
const (
	FulfillmentAwaiting       FulfillmentStatus = "AWAITING_PROCESSING"
	FulfillmentProcessing     FulfillmentStatus = "PROCESSING"
	FulfillmentShipped        FulfillmentStatus = "SHIPPED"
	FulfillmentDelivered      FulfillmentStatus = "DELIVERED"
	FulfillmentWillNotDeliver FulfillmentStatus = "WILL_NOT_DELIVER"
	FulfillmentReturned       FulfillmentStatus = "RETURNED"
	FulfillmentReadyForPickup FulfillmentStatus = "READY_FOR_PICKUP"
)

// OrderOptionType types
const (
	OptionChoice  OrderOptionType = "CHOICE"  // dropdown or radio button
	OptionChoices OrderOptionType = "CHOICES" // checkboxes
	OptionText    OrderOptionType = "TEXT"    // text input and text area
	OptionDate    OrderOptionType = "DATE"    // date/time
	OptionFiles   OrderOptionType = "FILES"   // upload file option
)

// SelectionModifierType types
const (
	SelectionModifierPercent  SelectionModifierType = "PERCENT"
	SelectionModifierAbsolute SelectionModifierType = "ABSOLUTE"
)
