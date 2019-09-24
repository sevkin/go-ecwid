package webhook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sevkin/go-ecwid"
)

// Event type
type Event string

// Event types
const (
	EventUnfinishedOrderCreated               Event = "unfinished_order.created"              // Unfinished order is created
	EventUnfinishedOrderUpdated               Event = "unfinished_order.updated"              // Unfinished order is updated
	EventUnfinishedOrderDeleted               Event = "unfinished_order.deleted"              // Unfinished order is deleted
	EventOrderCreated                         Event = "order.created"                         // New order is placed
	EventOrderUpdated                         Event = "order.updated"                         // Order is changed
	EventOrderDeleted                         Event = "order.deleted"                         // Order is deleted
	EventProductCreated                       Event = "product.created"                       // New product is created
	EventProductUpdated                       Event = "product.updated"                       // Product is updated
	EventProductDeleted                       Event = "product.deleted"                       // Product is deleted
	EventCategoryCreated                      Event = "category.created"                      // New category is created
	EventCategoryUpdated                      Event = "category.updated"                      // Category is updated
	EventCategoryDeleted                      Event = "category.deleted"                      // Category is deleted
	EventApplicationInstalled                 Event = "application.installed"                 // Application is installed
	EventApplicationUninstalled               Event = "application.uninstalled"               // Application is deleted
	EventApplicationSubscriptionStatusChanged Event = "application.subscriptionStatusChanged" // Application status changed
	EventProfileUpdated                       Event = "profile.updated"                       // Store information updated
	EventProfileSubscriptionStatusChanged     Event = "profile.subscriptionStatusChanged"     // Store premium subscription status changed
	EventCustomerCreated                      Event = "customer.created"                      // Customer is created
	EventCustomerUpdated                      Event = "customer.updated"                      // Customer is updated
	EventCustomerDeleted                      Event = "customer.deleted"                      // Customer is deleted
)

type (
	// Webhook like a http mux but ecwid webhook
	// https://developers.ecwid.com/api-documentation/webhook-structure
	Webhook interface {
		http.Handler
		Add(Event, Handler) Webhook // Add event handler
	}

	// Handler of webhook event
	Handler func(*Body) error

	// Body of webhook request
	Body struct {
		ID       string          `json:"eventId"`        // Unique webhook ID
		Event    Event           `json:"eventType"`      // Type of the occurred event
		Created  ecwid.Timestamp `json:"eventCreated"`   // Unix timestamp of the occurred event
		StoreID  ecwid.ID        `json:"storeId"`        // Store ID of the store where the event occured
		EntityID ecwid.ID        `json:"entityId"`       // Id of the updated entity. Can contain productId, categoryId, orderNumber, storeId depending on the eventType
		Data     *Data           `json:"data,omitempty"` // Describes changes made to an entity
	}

	// Data is Body Data optional field
	// Is provided for order.* and application.subscriptionStatusChanged event types
	Data struct {
		OldPaymentStatus      string `json:"oldPaymentStatus,omitempty"`      // Payment status of order before changes occurred
		NewPaymentStatus      string `json:"newPaymentStatus,omitempty"`      // Payment status of order after changes occurred
		OldFulfillmentStatus  string `json:"oldFulfillmentStatus,omitempty"`  // Fulfillment status of order before changes occurred
		NewFulfillmentStatus  string `json:"newFulfillmentStatus,omitempty"`  // Fulfillment status of an order after changes occurred
		OldSubscriptionName   string `json:"oldSubscriptionName,omitempty"`   // Previous Ecwid store premium plan name
		NewSubscriptionName   string `json:"newSubscriptionName,omitempty"`   // New Ecwid store premium plan name
		OldSubscriptionStatus string `json:"oldSubscriptionStatus,omitempty"` // Previous application subscription status before changes occurred
		NewSubscriptionStatus string `json:"newSubscriptionStatus,omitempty"` // New application subscription status after changes occurred
		CustomerEmail         string `json:"customerEmail,omitempty"`         // Email of a customer
	}

	webhook struct {
		secret string
		events map[Event]Handler
	}
)

// New returns new Webhook instance
// if clientSecret == "" don`t check secret
func New(clientSecret string) Webhook {
	return &webhook{
		secret: clientSecret,
		events: make(map[Event]Handler),
	}
}

func (wh *webhook) Add(event Event, handler Handler) Webhook {
	wh.events[event] = handler // TODO check if this event already exists
	return wh
}

func (wh *webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if wh.secret != "" {
		// TODO check X-Ecwid-Webhook-Signature
		// Encode the string ’{eventCreated}.{eventId}’
		// using HMAC SHA256 (using client_secret as the shared secret key)
		// and pass it through Base64 encoding
		//
		// see webhook_test.go TestOK
	}

	// TODO check Content-Type: application/json

	var body Body

	if buf, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal(buf, &body); err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if handler, found := wh.events[body.Event]; found {
		if err := handler(&body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
