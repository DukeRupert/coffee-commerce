package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID                uuid.UUID           `json:"id"`
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	ImageURL          string              `json:"image_url"`
	Active            bool                `json:"active"`
	StockLevel        int                 `json:"stock_level"`
	Weight            int                 `json:"weight"` // Base weight in grams
	Origin            string              `json:"origin"`
	RoastLevel        string              `json:"roast_level"`
	FlavorNotes       string              `json:"flavor_notes"`
	Options           map[string][]string `json:"options"`            // Product options (e.g., weight, grind)
	AllowSubscription bool                `json:"allow_subscription"` // Flag to indicate if product can be subscribed to
	StripeID          string              `json:"stripe_id"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

// WeightOption represents available weight options for coffee products
type WeightOption string

// GrindOption represents available grind options for coffee products
type GrindOption string

// Weight options constants
const (
	WeightTwelveOunce WeightOption = "12oz"
	WeightThreePound  WeightOption = "3lb"
	WeightFivePound   WeightOption = "5lb"
)

// Grind options constants
const (
	GrindWholeBeanOption GrindOption = "Whole Bean"
	GrindDripOption      GrindOption = "Drip Ground"
)

// Variant represents a specific product variant (combination of product, weight, and grind)
type Variant struct {
	ID            uuid.UUID `json:"id"`
	ProductID     uuid.UUID `json:"product_id"`
	PriceID       uuid.UUID `json:"price_id"`
	StripePriceID string    `json:"stripe_price_id"`
	Weight        string    `json:"weight"`
	Grind         string    `json:"grind"`
	Active        bool      `json:"active"`
	StockLevel    int       `json:"stock_level"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// VariantWithDetails includes product and price details for API responses
type VariantWithDetails struct {
	Variant
	ProductName  string `json:"product_name"`
	ProductImage string `json:"product_image,omitempty"`
	Origin       string `json:"origin,omitempty"`
	RoastLevel   string `json:"roast_level,omitempty"`
	FlavorNotes  string `json:"flavor_notes,omitempty"`
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
	PriceName    string `json:"price_name,omitempty"`
}

// GetWeightOptions returns all available weight options
func GetWeightOptions() []WeightOption {
	return []WeightOption{
		WeightTwelveOunce,
		WeightThreePound,
		WeightFivePound,
	}
}

// GetGrindOptions returns all available grind options
func GetGrindOptions() []GrindOption {
	return []GrindOption{
		GrindWholeBeanOption,
		GrindDripOption,
	}
}

// Price represents the pricing options for subscriptions or one-time purchases
type Price struct {
	ID            uuid.UUID `json:"id"`
	ProductID     uuid.UUID `json:"product_id"`
	Name          string    `json:"name"`
	Amount        int64     `json:"amount"` // Price in cents
	Currency      string    `json:"currency"`
	Type          string    `json:"type"`                     // one_time|recurring
	Interval      string    `json:"interval,omitempty"`       // week|month|year (used only for recurring)
	IntervalCount int       `json:"interval_count,omitempty"` // Number of intervals between charges (used only for recurring)
	Active        bool      `json:"active"`
	StripeID      string    `json:"stripe_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Customer represents a subscriber in the system
type Customer struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	StripeID    string    `json:"stripe_id"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Address represents a customer's shipping address
type Address struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Line1      string    `json:"line1"`
	Line2      string    `json:"line2"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

// Subscription status constants - matching Stripe's status values
const (
	SubscriptionStatusActive            = "active"
	SubscriptionStatusPastDue           = "past_due"
	SubscriptionStatusIncomplete        = "incomplete"
	SubscriptionStatusIncompleteExpired = "incomplete_expired"
	SubscriptionStatusTrialing          = "trialing"
	SubscriptionStatusCanceled          = "canceled" // Note the American spelling used by Stripe
	SubscriptionStatusUnpaid            = "unpaid"
	SubscriptionStatusPaused            = "paused" // Our custom status
)

// Subscription represents a customer's subscription to a product
type Subscription struct {
	ID         uuid.UUID  `json:"id"`
	CustomerID uuid.UUID  `json:"customer_id"`
	ProductID  uuid.UUID  `json:"product_id"`
	PriceID    uuid.UUID  `json:"price_id"`
	AddressID  *uuid.UUID `json:"address_id,omitempty"` // Optional for checkout process

	// Stripe subscription data
	StripeID     string `json:"stripe_id"`      // Main subscription ID
	StripeItemID string `json:"stripe_item_id"` // ID of the individual subscription item

	// Quantity and status
	Quantity int    `json:"quantity"`
	Status   string `json:"status"` // Using Stripe's status values

	// Billing period
	CurrentPeriodStart time.Time `json:"current_period_start"` // When the current billing period started
	CurrentPeriodEnd   time.Time `json:"current_period_end"`   // When the current billing period ends
	NextDeliveryDate   time.Time `json:"next_delivery_date"`   // When the coffee will be delivered

	// Cancellation details
	CancelAtPeriodEnd bool       `json:"cancel_at_period_end"`  // Whether to cancel at period end
	CanceledAt        *time.Time `json:"canceled_at,omitempty"` // When the subscription was canceled

	// Metadata
	Metadata  map[string]string `json:"metadata,omitempty"` // Additional data
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// SubscriptionWithDetails includes related entity details for API responses
type SubscriptionWithDetails struct {
	Subscription
	ProductName   string `json:"product_name"`
	ProductImage  string `json:"product_image,omitempty"`
	PriceName     string `json:"price_name"`
	Interval      string `json:"interval"`       // week, month, year
	IntervalCount int    `json:"interval_count"` // e.g., 2 for bi-weekly
	Amount        int64  `json:"amount"`         // Price amount in cents
	Currency      string `json:"currency"`       // USD, EUR, etc.
}
