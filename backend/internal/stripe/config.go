// internal/stripe/config.go
package stripe

import (
	"github.com/dukerupert/coffee-commerce/config"
	"github.com/rs/zerolog"
	stripeSDK "github.com/stripe/stripe-go/v82"
)

// InitStripe initializes the Stripe client with the provided API key
func InitStripe(cfg *config.StripeConfig, logger *zerolog.Logger) {
	subLogger := logger.With().Str("component", "stripe_init").Logger()
	
	if cfg.SecretKey == "" {
		subLogger.Warn().Msg("No Stripe secret key provided, Stripe functionality will be limited")
		return
	}
	
	// Set the API key for all API requests
	stripeSDK.Key = cfg.SecretKey
	
	subLogger.Info().Msg("Stripe SDK initialized successfully")
}