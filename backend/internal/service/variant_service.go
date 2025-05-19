// internal/service/variant_service.go
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	events "github.com/dukerupert/coffee-commerce/internal/event"
	interfaces "github.com/dukerupert/coffee-commerce/internal/repository/interface"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// VariantService defines the interface for variant-related operations
type VariantService interface {
	// Methods will be added later
}

// variantService implements VariantService
type variantService struct {
	logger      zerolog.Logger
	eventBus    events.EventBus
	variantRepo interfaces.VariantRepository
	productRepo interfaces.ProductRepository
	priceRepo   interfaces.PriceRepository
}

// NewVariantService creates a new variant service and subscribes to relevant events
func NewVariantService(logger *zerolog.Logger, eventBus events.EventBus, variantRepo interfaces.VariantRepository, productRepo interfaces.ProductRepository, priceRepo interfaces.PriceRepository) (VariantService, error) {
	subLogger := logger.With().Str("component", "variant_service").Logger()

	s := &variantService{
		logger:      subLogger,
		eventBus:    eventBus,
		variantRepo: variantRepo,
		productRepo: productRepo,
		priceRepo:   priceRepo,
	}

	// Subscribe to product created events
	_, err := eventBus.Subscribe(events.TopicProductCreated, s.handleProductCreated)
	if err != nil {
		subLogger.Error().Err(err).Msg("Failed to subscribe to product created events")
		return nil, err
	}

	// Subscribe to variant created events
	_, err = eventBus.Subscribe(events.TopicVariantQueued, s.handleVariantQueued)
	if err != nil {
		subLogger.Error().Err(err).Msg("Failed to subscribe to variant queued events")
		return nil, err
	}

	subLogger.Info().Str("topic", events.TopicProductCreated).Msg("Subscribed to product created events")
	return s, nil
}

// handleProductCreated is called when a product created event is received
func (s *variantService) handleProductCreated(data []byte) {
	s.logger.Info().Str("topic", events.TopicProductCreated).Msg("Received product created event")

	// Parse the event
	var event events.Event
	if err := json.Unmarshal(data, &event); err != nil {
		s.logger.Error().Err(err).Msg("Failed to unmarshal product created event")
		return
	}

	// Unmarshal the payload to get product details
	var payload events.ProductCreatedPayload
	payloadData, err := json.Marshal(event.Payload)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to marshal payload for unmarshaling")
		return
	}

	if err := json.Unmarshal(payloadData, &payload); err != nil {
		s.logger.Error().Err(err).Msg("Failed to unmarshal product created payload")
		return
	}

	// Log the received product details
	s.logger.Debug().
		Str("product_id", payload.ProductID).
		Str("name", payload.Name).
		Interface("options", payload.Options).
		Bool("allow_subscription", payload.AllowSubscription).
		Msg("Processing product created event")

	// Check if product has options
	if payload.Options == nil || len(payload.Options) == 0 {
		s.logger.Info().Str("product_id", payload.ProductID).
			Msg("Product has no options, creating single price")

		// Create a single price for the product via Stripe
		err := s.createDefaultPrice(payload.ProductID)
		if err != nil {
			s.logger.Error().Err(err).
				Str("product_id", payload.ProductID).
				Msg("Failed to create default price for product")
		}
		return
	}

	// Product has options, need to create variants for all combinations
	// Collect all option sets
	optionSets := make([][]string, 0, len(payload.Options))
	optionKeys := make([]string, 0, len(payload.Options))

	// Log all option keys and values
	for key, values := range payload.Options {
		if len(values) == 0 {
			s.logger.Warn().
				Str("product_id", payload.ProductID).
				Str("option_key", key).
				Msg("Option key has no values, skipping")
			continue
		}

		s.logger.Debug().
			Str("product_id", payload.ProductID).
			Str("option_key", key).
			Strs("option_values", values).
			Msg("Processing option set")

		optionSets = append(optionSets, values)
		optionKeys = append(optionKeys, key)
	}

	// Check if we have any valid option sets
	if len(optionSets) == 0 {
		s.logger.Info().
			Str("product_id", payload.ProductID).
			Msg("No valid options found, creating single price")

		// Create a single price for the product via Stripe
		err := s.createDefaultPrice(payload.ProductID)
		if err != nil {
			s.logger.Error().Err(err).
				Str("product_id", payload.ProductID).
				Msg("Failed to create default price for product")
		}
		return
	}

	// Generate all combinations of options
	combinations := s.generateOptionCombinations(optionSets)

	// Log the number of variants to be created
	s.logger.Info().
		Str("product_id", payload.ProductID).
		Int("option_sets", len(optionSets)).
		Int("total_combinations", len(combinations)).
		Msg("Found product options - will create variants for all combinations")

		// Save these combinations for later variant creation once prices are available
		// For now, just queue them for processing
	s.queueVariantCreation(payload.ProductID, optionKeys, combinations, payload)
}

// createDefaultPrice creates a default price for a product without options
func (s *variantService) createDefaultPrice(productID string) error {
	// TODO: Call Stripe API to create a price for the product
	// This is just a stub for now
	s.logger.Info().
		Str("product_id", productID).
		Msg("Stub: Would create a default price for product without options")
	return nil
}

// generateOptionCombinations generates all possible combinations of option values
func (s *variantService) generateOptionCombinations(optionSets [][]string) [][]string {
	if len(optionSets) == 0 {
		return [][]string{}
	}

	if len(optionSets) == 1 {
		// Base case: For a single option set, each value is its own combination
		result := make([][]string, len(optionSets[0]))
		for i, val := range optionSets[0] {
			result[i] = []string{val}
		}
		return result
	}

	// Get combinations of all but the first option set
	subCombinations := s.generateOptionCombinations(optionSets[1:])

	// Combine the first option set with all sub-combinations
	result := make([][]string, 0, len(optionSets[0])*len(subCombinations))

	for _, firstVal := range optionSets[0] {
		for _, subComb := range subCombinations {
			// Create a new combination with the current value from the first option set
			// followed by all values from the sub-combination
			newComb := make([]string, 1+len(subComb))
			newComb[0] = firstVal
			copy(newComb[1:], subComb)

			result = append(result, newComb)
		}
	}

	return result
}

// queueVariantCreation publishes events for each variant combination to be created
func (s *variantService) queueVariantCreation(productID string, optionKeys []string, combinations [][]string, payload events.ProductCreatedPayload) {
	s.logger.Info().
		Str("product_id", productID).
		Strs("option_keys", optionKeys).
		Int("combinations", len(combinations)).
		Msg("Publishing variant creation events to NATS")

	// Get a default price in cents (can be updated later)
	defaultPrice := int64(1000) // $10.00 by default
	defaultCurrency := "USD"

	// For each combination, create a payload and publish an event
	for i, combination := range combinations {
		// Convert the combination into a map of option key -> option value
		optionValues := make(map[string]string)
		for j, value := range combination {
			if j < len(optionKeys) {
				optionValues[optionKeys[j]] = value
			}
		}

		// Create a variant name that includes the options
		variantName := payload.Name
		for key, value := range optionValues {
			variantName += " - " + key + ": " + value
		}

		// Create the variant creation payload
		variantPayload := events.VariantQueuedPayload{
			ProductID:    productID,
			ProductName:  variantName,
			Description:  payload.Description,
			ImageURL:     payload.ImageURL,
			OptionValues: optionValues,
			DefaultPrice: defaultPrice,
			Currency:     defaultCurrency,
			QueuedAt:     time.Now(),
		}

		// Publish the event
		err := s.eventBus.Publish(events.TopicVariantQueued, variantPayload)
		if err != nil {
			s.logger.Error().Err(err).
				Str("product_id", productID).
				Int("combination_index", i).
				Interface("option_values", optionValues).
				Msg("Failed to publish variant queued event")
		} else {
			s.logger.Debug().
				Str("product_id", productID).
				Int("combination_index", i).
				Interface("option_values", optionValues).
				Msg("Published variant queued event")
		}
	}

	s.logger.Info().
		Str("product_id", productID).
		Int("total_variants", len(combinations)).
		Msg("Completed publishing variant creation events")
}

// handleVariantQueued processes events for variants that need Stripe products and prices
func (s *variantService) handleVariantQueued(data []byte) {
	// Parse the event
	var event events.Event
	if err := json.Unmarshal(data, &event); err != nil {
		s.logger.Error().Err(err).Msg("Failed to unmarshal variant queued event")
		return
	}

	// Unmarshal the payload
	var payload events.VariantQueuedPayload
	payloadData, err := json.Marshal(event.Payload)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to marshal payload for unmarshaling")
		return
	}

	if err := json.Unmarshal(payloadData, &payload); err != nil {
		s.logger.Error().Err(err).Msg("Failed to unmarshal variant queued payload")
		return
	}

	s.logger.Info().
		Str("product_id", payload.ProductID).
		Str("variant_name", payload.ProductName).
		Interface("option_values", payload.OptionValues).
		Msg("Processing variant creation with Stripe integration")

	// Create a Stripe product for this variant
	stripeProductObj, err := s.createStripeProduct(payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", payload.ProductID).
			Interface("option_values", payload.OptionValues).
			Msg("Failed to create Stripe product, will retry later")
		return
	}

	// Extract the stripe product ID from our mock object
	var stripeProductID string
	switch prod := stripeProductObj.(type) {
	case struct {
		ID          string
		Name        string
		Description string
		Metadata    map[string]string
	}:
		stripeProductID = prod.ID
	default:
		// In a real implementation with the Stripe SDK, you would use:
		// stripeProductID = stripeProductObj.ID
		stripeProductID = fmt.Sprintf("prod_%s", uuid.New().String())
	}

	// Create a default price for this product
	stripePriceObj, err := s.createStripePrice(stripeProductID, payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", payload.ProductID).
			Str("stripe_product_id", stripeProductID).
			Msg("Failed to create Stripe price, will retry later")
		return
	}

	// Extract the stripe price ID from our mock object
	var stripePriceID string
	switch price := stripePriceObj.(type) {
	case struct {
		ID         string
		ProductID  string
		UnitAmount int64
		Currency   string
	}:
		stripePriceID = price.ID
	default:
		// In a real implementation with the Stripe SDK, you would use:
		// stripePriceID = stripePriceObj.ID
		stripePriceID = fmt.Sprintf("price_%s", uuid.New().String())
	}

	// Create the variant in our database
	variant, err := s.createVariant(payload, stripeProductID, stripePriceID)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", payload.ProductID).
			Str("stripe_product_id", stripeProductID).
			Str("stripe_price_id", stripePriceID).
			Msg("Failed to create variant in database")
		return
	}

	// Publish an event that the variant was created
	variantCreatedPayload := events.VariantCreatedPayload{
		VariantID:     variant.ID.String(),
		ProductID:     payload.ProductID,
		StripeID:      stripeProductID,
		StripePriceID: stripePriceID,
		OptionValues:  payload.OptionValues,
		CreatedAt:     time.Now(),
	}

	err = s.eventBus.Publish(events.TopicVariantCreated, variantCreatedPayload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("variant_id", variant.ID.String()).
			Msg("Failed to publish variant created event")
	}

	s.logger.Info().
		Str("variant_id", variant.ID.String()).
		Str("product_id", payload.ProductID).
		Str("stripe_product_id", stripeProductID).
		Str("stripe_price_id", stripePriceID).
		Msg("Successfully created variant with Stripe integration")
}

// createStripeProduct creates a new product in Stripe
func (s *variantService) createStripeProduct(payload events.VariantQueuedPayload) (interface{}, error) {
	// Create the Stripe product
	// Note: This would typically use the Stripe Go SDK
	s.logger.Info().
		Str("product_name", payload.ProductName).
		Msg("Creating Stripe product")

	// In a real implementation, this would be:
	/*
	   params := &stripe.ProductParams{
	       Name:        stripe.String(payload.ProductName),
	       Description: stripe.String(payload.Description),
	   }

	   if payload.ImageURL != "" {
	       params.Images = []*string{stripe.String(payload.ImageURL)}
	   }

	   // Add option values as metadata
	   params.Metadata = make(map[string]string)
	   for key, value := range payload.OptionValues {
	       params.Metadata[key] = value
	   }

	   product, err := product.New(params)
	   if err != nil {
	       return nil, fmt.Errorf("failed to create Stripe product: %w", err)
	   }

	   return product, nil
	*/

	// For now, we'll simulate creating a Stripe product
	// Using a simple struct instead of the actual Stripe type
	mockProduct := struct {
		ID          string
		Name        string
		Description string
		Metadata    map[string]string
	}{
		ID:          "prod_" + uuid.New().String(),
		Name:        payload.ProductName,
		Description: payload.Description,
		Metadata:    payload.OptionValues,
	}

	return mockProduct, nil
}

// createStripePrice creates a new price in Stripe
func (s *variantService) createStripePrice(stripeProductID string, payload events.VariantQueuedPayload) (interface{}, error) {
	// Create the Stripe price
	s.logger.Info().
		Str("stripe_product_id", stripeProductID).
		Int64("price", payload.DefaultPrice).
		Str("currency", payload.Currency).
		Msg("Creating Stripe price")

	// In a real implementation, this would be:
	/*
	   params := &stripe.PriceParams{
	       Product:    stripe.String(stripeProductID),
	       UnitAmount: stripe.Int64(payload.DefaultPrice),
	       Currency:   stripe.String(payload.Currency),
	   }

	   price, err := price.New(params)
	   if err != nil {
	       return nil, fmt.Errorf("failed to create Stripe price: %w", err)
	   }

	   return price, nil
	*/

	// For now, we'll simulate creating a Stripe price
	// Using a simple struct instead of the actual Stripe type
	mockPrice := struct {
		ID         string
		ProductID  string
		UnitAmount int64
		Currency   string
	}{
		ID:         "price_" + uuid.New().String(),
		ProductID:  stripeProductID,
		UnitAmount: payload.DefaultPrice,
		Currency:   payload.Currency,
	}

	return mockPrice, nil
}

// createVariant creates a new variant in our database
func (s *variantService) createVariant(payload events.VariantQueuedPayload, stripeProductID string, stripePriceID string) (*model.Variant, error) {
	ctx := context.Background()

	// Ensure we have the required information
	if payload.ProductID == "" {
		return nil, errors.New("product ID is required")
	}

	// Parse product ID
	productID, err := uuid.Parse(payload.ProductID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	// Set sensible defaults for missing values
	price := int64(1000) // Default to $10.00
	if payload.DefaultPrice > 0 {
		price = payload.DefaultPrice
	}

	currency := "USD" // Default currency
	if payload.Currency != "" {
		currency = payload.Currency
	}

	// Create a new price record in our database
	priceRecord := &model.Price{
		ID:        uuid.New(),
		ProductID: productID,
		Name:      fmt.Sprintf("%s - %s", payload.ProductName, currency),
		Amount:    price,
		Currency:  currency,
		Type:      "one_time", // Default to one-time price
		Active:    true,
		StripeID:  stripePriceID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.logger.Debug().
		Str("price_id", priceRecord.ID.String()).
		Str("product_id", productID.String()).
		Int64("amount", price).
		Str("currency", currency).
		Msg("Creating price record")

	err = s.priceRepo.Create(ctx, priceRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to create price record: %w", err)
	}

	// Initialize the options map for the variant
	options := make(map[string]string)

	// Check if OptionValues is nil before accessing it
	if payload.OptionValues != nil {
		// Populate the options map with all options
		for key, value := range payload.OptionValues {
			options[key] = value
		}
	}

	// Create the variant without explicit weight and grind fields
	variant := &model.Variant{
		ID:            uuid.New(),
		ProductID:     productID,
		PriceID:       priceRecord.ID,
		StripePriceID: stripePriceID,
		Options:       options,
		Active:        true,
		StockLevel:    0, // Default to 0 until specifically set
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	s.logger.Debug().
		Str("variant_id", variant.ID.String()).
		Str("product_id", productID.String()).
		Str("price_id", priceRecord.ID.String()).
		Interface("options", options).
		Msg("Creating variant record")

	err = s.variantRepo.Create(ctx, variant)
	if err != nil {
		return nil, fmt.Errorf("failed to create variant: %w", err)
	}

	return variant, nil
}
