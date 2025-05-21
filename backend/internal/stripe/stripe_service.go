// internal/stripe/stripe_service.go
package stripe

import (
	"errors"
	"fmt"

	"github.com/dukerupert/coffee-commerce/config"
	"github.com/rs/zerolog"
	stripe "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
)

// Common errors
var (
	ErrStripeAPIDisabled = errors.New("stripe API is disabled or not properly configured")
)

// Service handles communication with the Stripe API
type Service struct {
	logger     zerolog.Logger
	config     *config.StripeConfig
	isDisabled bool
}

// NewStripeService creates a new Stripe service
func NewStripeService(logger *zerolog.Logger, cfg *config.StripeConfig) *Service {
	subLogger := logger.With().Str("component", "stripe_service").Logger()
	
	// Check if Stripe is properly configured
	isDisabled := cfg.SecretKey == ""
	if isDisabled {
		subLogger.Warn().Msg("Stripe service created in disabled mode - no API key provided")
	} else {
		// Initialize the Stripe SDK with the API key
		InitStripe(cfg, &subLogger)
	}
	
	return &Service{
		logger:     subLogger,
		config:     cfg,
		isDisabled: isDisabled,
	}
}

// CreateProduct creates a new product in Stripe
func (s *Service) CreateProduct(name, description string, imageURLs []string, metadata map[string]string) (*stripe.Product, error) {
    if s.isDisabled {
        s.logger.Warn().Msg("Stripe is disabled, returning mock product")
        return &stripe.Product{
            ID:          fmt.Sprintf("prod_mock_%s", name),
            Name:        name,
            Description: description,
            Metadata:    metadata,
        }, nil
    }
    
    s.logger.Debug().
        Str("name", name).
        Str("description", description).
        Strs("image_urls", imageURLs).
        Interface("metadata", metadata).
        Msg("Creating Stripe product")
        
    params := &stripe.ProductParams{
        Name:        stripe.String(name),
        Description: stripe.String(description),
    }
    
    // Add images if provided
    if len(imageURLs) > 0 {
        params.Images = make([]*string, len(imageURLs))
        for i, url := range imageURLs {
            params.Images[i] = stripe.String(url)
        }
    }
    
    // Add metadata if provided
    if len(metadata) > 0 {
        params.Metadata = make(map[string]string)
        for k, v := range metadata {
            params.Metadata[k] = v
        }
    }
    
    prod, err := product.New(params)
    if err != nil {
        s.logger.Error().Err(err).
            Str("name", name).
            Msg("Failed to create Stripe product")
        return nil, fmt.Errorf("failed to create Stripe product: %w", err)
    }
    
    s.logger.Info().
        Str("product_id", prod.ID).
        Str("name", prod.Name).
        Msg("Successfully created Stripe product")
        
    return prod, nil
}

// CreatePrice creates a new price in Stripe
func (s *Service) CreatePrice(productID string, unitAmount int64, currency string, recurring bool, 
    interval string, intervalCount int64) (*stripe.Price, error) {
    
    if s.isDisabled {
        s.logger.Warn().Msg("Stripe is disabled, returning mock price")
        return &stripe.Price{
            ID:        fmt.Sprintf("price_mock_%d_%s", unitAmount, currency),
            Product:   &stripe.Product{ID: productID},
            UnitAmount: unitAmount,
            Currency:  stripe.Currency(currency),
        }, nil
    }
    
    s.logger.Debug().
        Str("product_id", productID).
        Int64("unit_amount", unitAmount).
        Str("currency", currency).
        Bool("recurring", recurring).
        Str("interval", interval).
        Int64("interval_count", intervalCount).
        Msg("Creating Stripe price")
    
    params := &stripe.PriceParams{
        Product:    stripe.String(productID),
        UnitAmount: stripe.Int64(unitAmount),
        Currency:   stripe.String(currency),
    }
    
    // If it's a recurring price, set the recurring parameters
    if recurring {
        params.Recurring = &stripe.PriceRecurringParams{
            Interval:      stripe.String(interval),
            IntervalCount: stripe.Int64(intervalCount),
        }
    }
    
    p, err := price.New(params)
    if err != nil {
        s.logger.Error().Err(err).
            Str("product_id", productID).
            Int64("unit_amount", unitAmount).
            Msg("Failed to create Stripe price")
        return nil, fmt.Errorf("failed to create Stripe price: %w", err)
    }
    
    s.logger.Info().
        Str("price_id", p.ID).
        Str("product_id", p.Product.ID).
        Int64("unit_amount", p.UnitAmount).
        Msg("Successfully created Stripe price")
        
    return p, nil
}

// GetProduct retrieves a product from Stripe by ID
func (s *Service) GetProduct(productID string) (*stripe.Product, error) {
	if s.isDisabled {
		s.logger.Warn().Msg("Stripe is disabled, returning mock product")
		return &stripe.Product{
			ID: productID,
			Name: "Mock Product",
		}, nil
	}
	
	p, err := product.Get(productID, nil)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", productID).
			Msg("Failed to retrieve Stripe product")
		return nil, fmt.Errorf("failed to retrieve Stripe product: %w", err)
	}
	
	return p, nil
}