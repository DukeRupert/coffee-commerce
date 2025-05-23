// internal/events/payloads.go
package events

import (
	"time"
)

// ProductCreatedPayload represents the data in a product.created event
type ProductCreatedPayload struct {
	// Core identifiers
	ProductID string `json:"product_id"`

	// Base product information
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`

	// Product details
	StockLevel  int    `json:"stock_level"`
	Weight      int    `json:"weight"` // Base weight in grams
	Origin      string `json:"origin"`
	RoastLevel  string `json:"roast_level"`
	FlavorNotes string `json:"flavor_notes"`

	// Product configuration
	Options           map[string][]string `json:"options"` // Available options (weights, grinds)
	AllowSubscription bool                `json:"allow_subscription"`
	Active            bool                `json:"active"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
}


// ProductUpdatedPayload represents the data in a product.updated event
type ProductUpdatedPayload struct {
	// Core identifiers
	ProductID string `json:"product_id"`

	// Base product information
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`

	// Product details
	StockLevel  int    `json:"stock_level"`
	Weight      int    `json:"weight"` // Base weight in grams
	Origin      string `json:"origin"`
	RoastLevel  string `json:"roast_level"`
	FlavorNotes string `json:"flavor_notes"`

	// Product configuration
	Options           map[string][]string `json:"options"` // Available options (weights, grinds)
	AllowSubscription bool                `json:"allow_subscription"`
	Active            bool                `json:"active"`
	Archived          bool                `json:"archived"`

	// Update metadata
	UpdatedAt            time.Time `json:"updated_at"`
	OptionsChanged       bool      `json:"options_changed"`        // Whether options were modified
	ShouldCreateVariants bool      `json:"should_create_variants"` // Whether this update should trigger variant creation
}

// ProductStockUpdatedPayload represents the data in a product.stock_updated event
type ProductStockUpdatedPayload struct {
	ProductID     string    `json:"product_id"`
	Name          string    `json:"name"`
	OldStockLevel int       `json:"old_stock_level"`
	NewStockLevel int       `json:"new_stock_level"`
	IsLowStock    bool      `json:"is_low_stock"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// VariantQueuedPayload represents the data needed to create a variant
type VariantQueuedPayload struct {
	// IDs
	ProductID string `json:"product_id"`

	// Product base information (to help create meaningful Stripe products)
	ProductName string `json:"product_name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`

	// The specific variant configuration
	OptionValues map[string]string `json:"option_values"`

	// Price information (can be updated later)
	DefaultPrice int64  `json:"default_price"` // Default price in cents
	Currency     string `json:"currency"`      // Default: USD

	// Metadata
	QueuedAt time.Time `json:"queued_at"`
}

// VariantCreatedPayload represents the data in a variant.created event
type VariantCreatedPayload struct {
	// IDs
	VariantID string `json:"variant_id"`
	ProductID string `json:"product_id"`
	PriceID   string `json:"price_id"`

	// Stripe IDs
	StripeProductID      string `json:"stripe_id"` // Stripe product ID
	StripePriceID string `json:"stripe_price_id"`

	// Variant details
	Weight       string            `json:"weight"`
	Grind        string            `json:"grind"`
	OptionValues map[string]string `json:"option_values"` // All option key-value pairs

	// Pricing information
	Amount   int64  `json:"amount"`   // Price amount in cents
	Currency string `json:"currency"` // Currency code (e.g., USD)

	// Status
	Active     bool `json:"active"`
	StockLevel int  `json:"stock_level"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
}

// VariantUpdatedPayload represents the data in a variant.updated event
type VariantUpdatedPayload struct {
	// IDs
	VariantID       string `json:"variant_id"`
	ProductID       string `json:"product_id"`
	PriceID         string `json:"price_id"`

	// Stripe IDs
	StripeProductID string `json:"stripe_product_id"`
	StripePriceID   string `json:"stripe_price_id"`

	// Price information
	Amount        int64  `json:"amount"`         // Price amount in cents
	Currency      string `json:"currency"`       // Currency code (e.g., USD)
	PriceType     string `json:"price_type"`     // one_time or recurring
	Interval      string `json:"interval,omitempty"`       // week, month, year (for recurring)
	IntervalCount int    `json:"interval_count,omitempty"` // Number of intervals (for recurring)

	// Metadata
	UpdatedAt    time.Time `json:"updated_at"`
	UpdateSource string    `json:"update_source"` // e.g., "stripe_webhook", "api", "admin"
}
