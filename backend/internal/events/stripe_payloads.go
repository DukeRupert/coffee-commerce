// internal/event/stripe_payloads.go
package events

import (
	"time"
)

// StripeProductEventPayload represents the data structure for Stripe product events
type StripeProductEventPayload struct {
	StripeID    string            `json:"stripe_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Active      bool              `json:"active"`
	Images      []string          `json:"images"`
	Metadata    map[string]string `json:"metadata"`
	URL         string            `json:"url"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// StripePriceEventPayload represents the data structure for Stripe price events
type StripePriceEventPayload struct {
	StripeID    string            `json:"stripe_id"`
	ProductID   string            `json:"product_id"`
	UnitAmount  int64             `json:"unit_amount"`
	Currency    string            `json:"currency"`
	Active      bool              `json:"active"`
	Type        string            `json:"type"` // one_time or recurring
	Recurring   *StripeRecurring  `json:"recurring,omitempty"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// StripeRecurring represents recurring pricing information
type StripeRecurring struct {
	Interval      string `json:"interval"`       // day, week, month, year
	IntervalCount int    `json:"interval_count"` // number of intervals
}

// StripeCustomerEventPayload represents the data structure for Stripe customer events
type StripeCustomerEventPayload struct {
	StripeID    string            `json:"stripe_id"`
	Email       string            `json:"email"`
	Name        string            `json:"name"`
	Phone       string            `json:"phone"`
	Address     *StripeAddress    `json:"address,omitempty"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// StripeAddress represents a customer address
type StripeAddress struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// StripeSubscriptionEventPayload represents the data structure for Stripe subscription events
type StripeSubscriptionEventPayload struct {
	StripeID           string                 `json:"stripe_id"`
	CustomerID         string                 `json:"customer_id"`
	Status             string                 `json:"status"`
	CurrentPeriodStart time.Time              `json:"current_period_start"`
	CurrentPeriodEnd   time.Time              `json:"current_period_end"`
	CancelAtPeriodEnd  bool                   `json:"cancel_at_period_end"`
	CanceledAt         *time.Time             `json:"canceled_at,omitempty"`
	Items              []StripeSubscriptionItem `json:"items"`
	Metadata           map[string]string      `json:"metadata"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

// StripeSubscriptionItem represents an item in a subscription
type StripeSubscriptionItem struct {
	StripeID  string `json:"stripe_id"`
	PriceID   string `json:"price_id"`
	Quantity  int64  `json:"quantity"`
	ProductID string `json:"product_id"`
}

// StripeCheckoutEventPayload represents the data structure for Stripe checkout events
type StripeCheckoutEventPayload struct {
	StripeID        string            `json:"stripe_id"`
	CustomerID      string            `json:"customer_id,omitempty"`
	CustomerEmail   string            `json:"customer_email,omitempty"`
	CustomerDetails *StripeCustomerDetails `json:"customer_details,omitempty"`
	PaymentStatus   string            `json:"payment_status"`
	Mode            string            `json:"mode"` // payment, setup, subscription
	LineItems       []StripeLineItem  `json:"line_items,omitempty"`
	Metadata        map[string]string `json:"metadata"`
	CreatedAt       time.Time         `json:"created_at"`
}

// StripeCustomerDetails contains customer information from a checkout session
type StripeCustomerDetails struct {
	Email   string        `json:"email"`
	Name    string        `json:"name,omitempty"`
	Phone   string        `json:"phone,omitempty"`
	Address *StripeAddress `json:"address,omitempty"`
}

// StripeLineItem represents a line item in a checkout session
type StripeLineItem struct {
	PriceID   string `json:"price_id"`
	ProductID string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
}

// StripeInvoiceEventPayload represents the data structure for Stripe invoice events
type StripeInvoiceEventPayload struct {
	StripeID      string    `json:"stripe_id"`
	CustomerID    string    `json:"customer_id"`
	SubscriptionID string   `json:"subscription_id,omitempty"`
	Status        string    `json:"status"`
	AmountDue     int64     `json:"amount_due"`
	AmountPaid    int64     `json:"amount_paid"`
	Currency      string    `json:"currency"`
	Paid          bool      `json:"paid"`
	Attempted     bool      `json:"attempted"`
	Created       time.Time `json:"created"`
}