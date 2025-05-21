// internal/service/variant_service.go
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	events "github.com/dukerupert/coffee-commerce/internal/event"
	"github.com/dukerupert/coffee-commerce/internal/interfaces"
	"github.com/dukerupert/coffee-commerce/internal/stripe"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	stripeSDK "github.com/stripe/stripe-go/v82"
)

// variantService implements VariantService
type variantService struct {
	logger        zerolog.Logger
	eventBus      events.EventBus
	variantRepo   interfaces.VariantRepository
	productRepo   interfaces.ProductRepository
	priceRepo     interfaces.PriceRepository
	stripeService *stripe.Service
}

// NewVariantService creates a new variant service and subscribes to relevant events
func NewVariantService(logger *zerolog.Logger, eventBus events.EventBus, variantRepo interfaces.VariantRepository, productRepo interfaces.ProductRepository, priceRepo interfaces.PriceRepository, stripeService *stripe.Service) (interfaces.VariantService, error) {
	subLogger := logger.With().Str("component", "variant_service").Logger()

	s := &variantService{
		logger:        subLogger,
		eventBus:      eventBus,
		variantRepo:   variantRepo,
		productRepo:   productRepo,
		priceRepo:     priceRepo,
		stripeService: stripeService,
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

// createDefaultPrice creates a default price and variant for a product without options
func (s *variantService) createDefaultPrice(productID string) error {
	ctx := context.Background()

	// Parse product ID
	prodID, err := uuid.Parse(productID)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	// Get product details
	product, err := s.productRepo.GetByID(ctx, prodID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return fmt.Errorf("product not found: %s", productID)
	}

	// Create a default price in our database
	priceID := uuid.New()
	stripePriceID := fmt.Sprintf("price_%s", uuid.New().String())

	price := &model.Price{
		ID:            priceID,
		ProductID:     prodID,
		Name:          fmt.Sprintf("%s - Default", product.Name),
		Amount:        1000, // Default $10.00
		Currency:      "USD",
		Type:          "one_time", // Default to one-time price
		Interval:      "",
		IntervalCount: 0,
		Active:        true,
		StripeID:      stripePriceID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	s.logger.Debug().
		Str("price_id", price.ID.String()).
		Str("product_id", prodID.String()).
		Int64("amount", price.Amount).
		Str("currency", price.Currency).
		Msg("Creating default price record")

	err = s.priceRepo.Create(ctx, price)
	if err != nil {
		return fmt.Errorf("failed to create price record: %w", err)
	}

	// Create a default variant
	variant := &model.Variant{
		ID:            uuid.New(),
		ProductID:     prodID,
		PriceID:       priceID,
		StripePriceID: stripePriceID,
		Weight:        product.Weight,
		Options:       make(map[string]string), // Empty options for default variant
		Active:        true,
		StockLevel:    product.StockLevel, // Use product's initial stock level
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	s.logger.Debug().
		Str("variant_id", variant.ID.String()).
		Str("product_id", prodID.String()).
		Str("price_id", priceID.String()).
		Msg("Creating default variant record")

	err = s.variantRepo.Create(ctx, variant)
	if err != nil {
		return fmt.Errorf("failed to create default variant: %w", err)
	}

	// Publish variant created event
	variantCreatedPayload := events.VariantCreatedPayload{
		VariantID:     variant.ID.String(),
		ProductID:     productID,
		PriceID:       priceID.String(),
		StripeID:      "prod_" + uuid.New().String(), // Simulate Stripe product ID
		StripePriceID: stripePriceID,
		Weight:        fmt.Sprintf("%dg", product.Weight),
		Grind:         "Whole Bean", // Default grind type
		OptionValues:  map[string]string{},
		Amount:        price.Amount,
		Currency:      price.Currency,
		Active:        true,
		StockLevel:    product.StockLevel,
		CreatedAt:     time.Now(),
	}

	err = s.eventBus.Publish(events.TopicVariantCreated, variantCreatedPayload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("variant_id", variant.ID.String()).
			Msg("Failed to publish default variant created event")
		// Don't return the error since the variant is already created in DB
	}

	s.logger.Info().
		Str("variant_id", variant.ID.String()).
		Str("product_id", productID).
		Str("price_id", priceID.String()).
		Msg("Successfully created default variant")

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
	stripeProduct, err := s.createStripeProduct(payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", payload.ProductID).
			Interface("option_values", payload.OptionValues).
			Msg("Failed to create Stripe product, will retry later")
		return
	}

	// Extract the Stripe product ID directly
	stripeProductID := stripeProduct.ID

	// Create a default price for this product
	stripePrice, err := s.createStripePrice(stripeProductID, payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", payload.ProductID).
			Str("stripe_product_id", stripeProductID).
			Msg("Failed to create Stripe price, will retry later")
		return
	}

	// Extract the Stripe price ID directly
	stripePriceID := stripePrice.ID

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
		PriceID:       variant.PriceID.String(),  // Include the internal price ID
		StripeID:      stripeProductID,
		StripePriceID: stripePriceID,
		Weight:        getOptionValue(payload.OptionValues, "weight", ""),
		Grind:         getOptionValue(payload.OptionValues, "grind", ""),
		OptionValues:  payload.OptionValues,
		Amount:        stripePrice.UnitAmount,
		Currency:      string(stripePrice.Currency),
		Active:        variant.Active,
		StockLevel:    variant.StockLevel,
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
		Int64("price_amount", stripePrice.UnitAmount).
		Str("currency", string(stripePrice.Currency)).
		Msg("Successfully created variant with Stripe integration")
}

// Helper function to safely get option values with a default fallback
func getOptionValue(options map[string]string, key, defaultValue string) string {
	if value, exists := options[key]; exists {
		return value
	}
	return defaultValue
}

// createStripeProduct creates a new product in Stripe
func (s *variantService) createStripeProduct(payload events.VariantQueuedPayload) (*stripeSDK.Product, error) {
	// Prepare metadata from options
	metadata := make(map[string]string)
	for k, v := range payload.OptionValues {
		metadata[k] = v
	}
	
	// Add original product ID to metadata
	metadata["original_product_id"] = payload.ProductID
	
	// Create images array if image URL is provided
	var images []string
	if payload.ImageURL != "" {
		images = []string{payload.ImageURL}
	}
	
	// Create the product in Stripe
	product, err := s.stripeService.CreateProduct(
		payload.ProductName,
		payload.Description,
		images,
		metadata,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe product: %w", err)
	}
	
	return product, nil
}

// createStripePrice creates a new price in Stripe
func (s *variantService) createStripePrice(stripeProductID string, payload events.VariantQueuedPayload) (*stripeSDK.Price, error) {
	// Determine if this is a recurring price
	recurring := false
	interval := ""
	intervalCount := int64(0)
	
	// Check for subscription options in the variant's metadata
	// If this is a coffee subscription, it would typically be monthly or bi-weekly
	if value, exists := payload.OptionValues["subscription_interval"]; exists {
		recurring = true
		
		// Parse interval from metadata - should be 'day', 'week', 'month', or 'year'
		interval = value
		
		// Get interval count if available, default to 1
		if countValue, exists := payload.OptionValues["subscription_interval_count"]; exists {
			if count, err := strconv.ParseInt(countValue, 10, 64); err == nil {
				intervalCount = count
			} else {
				intervalCount = 1
			}
		} else {
			intervalCount = 1
		}
	}
	
	// Set a sane default price if none provided
	amount := payload.DefaultPrice
	if amount <= 0 {
		// Default to $10.00 if no price specified
		amount = 1000
		s.logger.Warn().
			Str("product_id", payload.ProductID).
			Msg("No price specified for variant, defaulting to $10.00")
	}
	
	// Set default currency if none provided
	currency := payload.Currency
	if currency == "" {
		currency = "USD"
	}
	
	// Create the price in Stripe
	price, err := s.stripeService.CreatePrice(
		stripeProductID,
		amount,
		currency,
		recurring,
		interval,
		intervalCount,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe price: %w", err)
	}
	
	return price, nil
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

	// Check if this is a subscription price
	if value, exists := payload.OptionValues["subscription_interval"]; exists {
		priceRecord.Type = "recurring"
		priceRecord.Interval = value
		
		// Get interval count if available
		if countValue, exists := payload.OptionValues["subscription_interval_count"]; exists {
			if count, err := strconv.Atoi(countValue); err == nil {
				priceRecord.IntervalCount = count
			} else {
				priceRecord.IntervalCount = 1 
			}
		} else {
			priceRecord.IntervalCount = 1
		}
	}

	s.logger.Debug().
		Str("price_id", priceRecord.ID.String()).
		Str("product_id", productID.String()).
		Int64("amount", price).
		Str("currency", currency).
		Str("type", priceRecord.Type).
		Str("interval", priceRecord.Interval).
		Int("interval_count", priceRecord.IntervalCount).
		Msg("Creating price record")

	err = s.priceRepo.Create(ctx, priceRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to create price record: %w", err)
	}

	// Initialize the options map for the variant
	options := make(map[string]string)

	// Default weight in grams
	weight := 1 

	// Check if OptionValues is nil before accessing it
	if payload.OptionValues != nil {
		// Populate the options map with all options
		for key, value := range payload.OptionValues {
			options[key] = value
		}
		
		// Try to parse weight from options if present
		if weightStr, ok := options["weight"]; ok {
			// Convert weight option (like "12oz") to grams if possible
			weight = convertWeightToGrams(weightStr)
		}
	}

	// Create the variant record
	variant := &model.Variant{
		ID:            uuid.New(),
		ProductID:     productID,
		PriceID:       priceRecord.ID,
		StripePriceID: stripePriceID,
		Weight:        weight,
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
		Str("stripe_price_id", stripePriceID).
		Int("weight", weight).
		Interface("options", options).
		Msg("Creating variant record")

	err = s.variantRepo.Create(ctx, variant)
	if err != nil {
		return nil, fmt.Errorf("failed to create variant: %w", err)
	}

	return variant, nil
}

// convertWeightToGrams converts weight strings like "12oz" to grams
func convertWeightToGrams(weightStr string) int {
	// Common conversion factors
	const (
		gramsPerOz  = 28
		gramsPerLb  = 454
	)
	
	// Try to parse simple cases like "12oz" or "3lb"
	if strings.HasSuffix(weightStr, "oz") {
		value := strings.TrimSuffix(weightStr, "oz")
		if oz, err := strconv.Atoi(value); err == nil {
			return oz * gramsPerOz
		}
	}
	
	if strings.HasSuffix(weightStr, "lb") {
		value := strings.TrimSuffix(weightStr, "lb")
		if lb, err := strconv.Atoi(value); err == nil {
			return lb * gramsPerLb
		}
	}
	
	// Default to 1 gram if we can't parse
	return 1
}
