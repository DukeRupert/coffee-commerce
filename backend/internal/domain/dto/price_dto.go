// internal/domain/dto/price_dto.go
package dto

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

// Valid price types
var validPriceTypes = map[string]bool{
	"one_time":  true,
	"recurring": true,
}

// Valid recurring intervals
var validIntervals = map[string]bool{
	"week":  true,
	"month": true,
	"year":  true,
}

// PriceCreateDTO represents the data needed to create a new price
type PriceCreateDTO struct {
	ProductID     uuid.UUID `json:"product_id"`
	Name          string    `json:"name"`
	Amount        int64     `json:"amount"`                   // Price in cents
	Currency      string    `json:"currency"`                 // Default: USD
	Type          string    `json:"type"`                     // one_time or recurring
	Interval      string    `json:"interval,omitempty"`       // week, month, year (for recurring only)
	IntervalCount int       `json:"interval_count,omitempty"` // Number of intervals (for recurring only)
	Active        bool      `json:"active"`
}

// Valid validates the PriceCreateDTO
func (p *PriceCreateDTO) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	// Validate required fields
	if p.ProductID == uuid.Nil {
		problems["product_id"] = "product ID is required"
	}

	if p.Name == "" {
		problems["name"] = "name is required"
	} else if len(p.Name) > 255 {
		problems["name"] = "name must not exceed 255 characters"
	}

	if p.Amount <= 0 {
		problems["amount"] = "amount must be greater than 0"
	} else if p.Amount > 99999999 { // $999,999.99 max
		problems["amount"] = "amount must not exceed $999,999.99"
	}

	// Validate currency
	if p.Currency == "" {
		p.Currency = "USD" // Default currency
	} else {
		p.Currency = strings.ToUpper(p.Currency)
		if len(p.Currency) != 3 {
			problems["currency"] = "currency must be a 3-letter code (e.g., USD, EUR)"
		}
	}

	// Validate price type
	if p.Type == "" {
		problems["type"] = "price type is required"
	} else if !validPriceTypes[strings.ToLower(p.Type)] {
		problems["type"] = "price type must be 'one_time' or 'recurring'"
	}

	// Validate recurring-specific fields
	if strings.ToLower(p.Type) == "recurring" {
		if p.Interval == "" {
			problems["interval"] = "interval is required for recurring prices"
		} else if !validIntervals[strings.ToLower(p.Interval)] {
			problems["interval"] = "interval must be 'week', 'month', or 'year'"
		}

		if p.IntervalCount <= 0 {
			p.IntervalCount = 1 // Default to 1
		} else if p.IntervalCount > 12 {
			problems["interval_count"] = "interval count must not exceed 12"
		}
	} else {
		// For one-time prices, these fields should be empty
		if p.Interval != "" {
			problems["interval"] = "interval should not be set for one-time prices"
		}
		if p.IntervalCount > 0 {
			problems["interval_count"] = "interval count should not be set for one-time prices"
		}
	}

	return problems
}

// ToModel converts PriceCreateDTO to a Price model
func (p *PriceCreateDTO) ToModel() *model.Price {
	price := &model.Price{
		ID:        uuid.New(),
		ProductID: p.ProductID,
		Name:      p.Name,
		Amount:    p.Amount,
		Currency:  strings.ToUpper(p.Currency),
		Type:      strings.ToLower(p.Type),
		Active:    p.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set recurring fields only if type is recurring
	if price.Type == "recurring" {
		price.Interval = strings.ToLower(p.Interval)
		if p.IntervalCount <= 0 {
			price.IntervalCount = 1
		} else {
			price.IntervalCount = p.IntervalCount
		}
	}

	return price
}

// PriceUpdateDTO represents the data needed to update an existing price
type PriceUpdateDTO struct {
	Name          *string `json:"name,omitempty"`
	Amount        *int64  `json:"amount,omitempty"`
	Currency      *string `json:"currency,omitempty"`
	Type          *string `json:"type,omitempty"`
	Interval      *string `json:"interval,omitempty"`
	IntervalCount *int    `json:"interval_count,omitempty"`
	Active        *bool   `json:"active,omitempty"`
}

// Valid validates the PriceUpdateDTO
func (p *PriceUpdateDTO) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	// Validate name if provided
	if p.Name != nil {
		if *p.Name == "" {
			problems["name"] = "name cannot be empty when provided"
		} else if len(*p.Name) > 255 {
			problems["name"] = "name must not exceed 255 characters"
		}
	}

	// Validate amount if provided
	if p.Amount != nil {
		if *p.Amount <= 0 {
			problems["amount"] = "amount must be greater than 0"
		} else if *p.Amount > 99999999 {
			problems["amount"] = "amount must not exceed $999,999.99"
		}
	}

	// Validate currency if provided
	if p.Currency != nil {
		currency := strings.ToUpper(*p.Currency)
		if len(currency) != 3 {
			problems["currency"] = "currency must be a 3-letter code (e.g., USD, EUR)"
		}
		*p.Currency = currency
	}

	// Validate price type if provided
	if p.Type != nil {
		priceType := strings.ToLower(*p.Type)
		if !validPriceTypes[priceType] {
			problems["type"] = "price type must be 'one_time' or 'recurring'"
		}
		*p.Type = priceType
	}

	// Validate interval if provided
	if p.Interval != nil {
		interval := strings.ToLower(*p.Interval)
		if !validIntervals[interval] {
			problems["interval"] = "interval must be 'week', 'month', or 'year'"
		}
		*p.Interval = interval
	}

	// Validate interval count if provided
	if p.IntervalCount != nil {
		if *p.IntervalCount <= 0 {
			problems["interval_count"] = "interval count must be greater than 0"
		} else if *p.IntervalCount > 12 {
			problems["interval_count"] = "interval count must not exceed 12"
		}
	}

	return problems
}

// ApplyToModel applies the non-nil fields from the DTO to the price model
func (p *PriceUpdateDTO) ApplyToModel(price *model.Price) {
	if p.Name != nil {
		price.Name = *p.Name
	}
	if p.Amount != nil {
		price.Amount = *p.Amount
	}
	if p.Currency != nil {
		price.Currency = *p.Currency
	}
	if p.Type != nil {
		price.Type = *p.Type
		// If changing to one_time, clear recurring fields
		if price.Type == "one_time" {
			price.Interval = ""
			price.IntervalCount = 0
		}
	}
	if p.Interval != nil {
		price.Interval = *p.Interval
	}
	if p.IntervalCount != nil {
		price.IntervalCount = *p.IntervalCount
	}
	if p.Active != nil {
		price.Active = *p.Active
	}
	price.UpdatedAt = time.Now()
}

// PriceResponseDTO represents the data returned to the client
type PriceResponseDTO struct {
	ID            string `json:"id"`
	ProductID     string `json:"product_id"`
	Name          string `json:"name"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Type          string `json:"type"`
	Interval      string `json:"interval,omitempty"`
	IntervalCount int    `json:"interval_count,omitempty"`
	Active        bool   `json:"active"`
	StripeID      string `json:"stripe_id,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`

	// Computed fields for convenience
	FormattedAmount string `json:"formatted_amount"` // e.g., "$10.00"
	DisplayName     string `json:"display_name"`     // e.g., "Monthly Subscription - $15.00/month"
	IsSubscription  bool   `json:"is_subscription"`  // true if type is recurring
}

// FromModel converts a Price model to PriceResponseDTO
func PriceResponseDTOFromModel(price *model.Price) PriceResponseDTO {
	response := PriceResponseDTO{
		ID:             price.ID.String(),
		ProductID:      price.ProductID.String(),
		Name:           price.Name,
		Amount:         price.Amount,
		Currency:       price.Currency,
		Type:           price.Type,
		Interval:       price.Interval,
		IntervalCount:  price.IntervalCount,
		Active:         price.Active,
		StripeID:       price.StripeID,
		CreatedAt:      price.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      price.UpdatedAt.Format(time.RFC3339),
		IsSubscription: price.Type == "recurring",
	}

	// Format amount for display
	response.FormattedAmount = formatCurrency(price.Amount, price.Currency)

	// Create display name
	response.DisplayName = createDisplayName(price)

	return response
}

// Helper function to format currency amounts
func formatCurrency(amountCents int64, currency string) string {
	amount := float64(amountCents) / 100.0

	switch strings.ToUpper(currency) {
	case "USD":
		return fmt.Sprintf("$%.2f", amount)
	case "EUR":
		return fmt.Sprintf("€%.2f", amount)
	case "GBP":
		return fmt.Sprintf("£%.2f", amount)
	default:
		return fmt.Sprintf("%.2f %s", amount, strings.ToUpper(currency))
	}
}

// Helper function to create a display name for the price
func createDisplayName(price *model.Price) string {
	formattedAmount := formatCurrency(price.Amount, price.Currency)

	if price.Type == "recurring" {
		interval := price.Interval
		if price.IntervalCount > 1 {
			// Handle pluralization
			switch interval {
			case "week":
				interval = fmt.Sprintf("%d weeks", price.IntervalCount)
			case "month":
				interval = fmt.Sprintf("%d months", price.IntervalCount)
			case "year":
				interval = fmt.Sprintf("%d years", price.IntervalCount)
			}
		} else {
			// Handle singular/special cases
			switch interval {
			case "month":
				interval = "month"
			case "week":
				interval = "week"
			case "year":
				interval = "year"
			}
		}

		return fmt.Sprintf("%s - %s/%s", price.Name, formattedAmount, interval)
	}

	return fmt.Sprintf("%s - %s (one-time)", price.Name, formattedAmount)
}

// PriceListResponseDTO represents a list of prices with pagination
type PriceListResponseDTO struct {
	Data []*PriceResponseDTO `json:"data"`
	Meta interface{}         `json:"meta,omitempty"` // Pagination metadata
}

// VariantPriceAssignmentDTO represents the data needed to assign a price to a variant
type VariantPriceAssignmentDTO struct {
	VariantID uuid.UUID `json:"variant_id"`
	PriceID   uuid.UUID `json:"price_id"`
}

// Valid validates the VariantPriceAssignmentDTO
func (v *VariantPriceAssignmentDTO) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if v.VariantID == uuid.Nil {
		problems["variant_id"] = "variant ID is required"
	}

	if v.PriceID == uuid.Nil {
		problems["price_id"] = "price ID is required"
	}

	return problems
}
