// internal/api/handler/stripe_webhook_handler.go
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dukerupert/coffee-commerce/config"
	"github.com/dukerupert/coffee-commerce/internal/api"
	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	events "github.com/dukerupert/coffee-commerce/internal/event"
	"github.com/dukerupert/coffee-commerce/internal/interfaces"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripeWebhookHandler struct {
	logger       zerolog.Logger
	stripeConfig *config.StripeConfig
	eventBus     events.EventBus
	productRepo  interfaces.ProductRepository
	priceRepo    interfaces.PriceRepository
	variantRepo  interfaces.VariantRepository
}

func NewStripeWebhookHandler(
	logger *zerolog.Logger,
	stripeConfig *config.StripeConfig,
	eventBus events.EventBus, productRepo interfaces.ProductRepository, priceRepo interfaces.PriceRepository,
	variantRepo interfaces.VariantRepository) *StripeWebhookHandler {

	return &StripeWebhookHandler{
		logger:       logger.With().Str("component", "stripe_webhook_handler").Logger(),
		stripeConfig: stripeConfig,
		eventBus:     eventBus,
		productRepo:  productRepo,
		priceRepo:    priceRepo,
		variantRepo:  variantRepo,
	}
}

// HandleWebhook handles Stripe webhook events
func (h *StripeWebhookHandler) HandleWebhook(c echo.Context) error {
	// Set a reasonable body size limit to prevent abuse
	const MaxBodyBytes = int64(65536)
	c.Request().Body = http.MaxBytesReader(c.Response().Writer, c.Request().Body, MaxBodyBytes)

	// Read the request body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to read webhook request body")
		return c.JSON(http.StatusServiceUnavailable, api.ErrorResponse{
			Status:  http.StatusServiceUnavailable,
			Message: "Failed to read request body",
		})
	}

	// Get the signature header
	signatureHeader := c.Request().Header.Get("Stripe-Signature")
	if signatureHeader == "" {
		h.logger.Warn().Msg("Missing Stripe-Signature header")
		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Missing Stripe-Signature header",
		})
	}

	// Verify the signature
	event, err := webhook.ConstructEvent(body, signatureHeader, h.stripeConfig.WebhookSecret)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to verify webhook signature")
		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to verify webhook signature",
		})
	}

	// Log the event
	h.logger.Info().
		Str("event_id", event.ID).
		Str("event_type", string(event.Type)).
		Msg("Received Stripe webhook event")

	// Process the event based on its type
	err = h.processEvent(event)
	if err != nil {
		h.logger.Error().Err(err).
			Str("event_id", event.ID).
			Str("event_type", string(event.Type)).
			Msg("Error processing webhook event")

		// Return 200 even if processing fails to prevent Stripe from retrying
		// We'll handle the error internally through logging and monitoring
	}

	// Always return a 200 OK to prevent Stripe from retrying
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

// processEvent handles different Stripe event types
func (h *StripeWebhookHandler) processEvent(event stripe.Event) error {
	switch event.Type {
	// Checkout session events
	case "checkout.session.async_payment_failed":
		return h.handleCheckoutSessionAsyncPaymentFailed(event)
	case "checkout.session.async_payment_succeeded":
		return h.handleCheckoutSessionAsyncPaymentSucceeded(event)
	case "checkout.session.completed":
		return h.handleCheckoutSessionCompleted(event)
	case "checkout.session.expired":
		return h.handleCheckoutSessionExpired(event)

	// Person events
	case "person.created":
		return h.handlePersonCreated(event)
	case "person.deleted":
		return h.handlePersonDeleted(event)
	case "person.updated":
		return h.handlePersonUpdated(event)

	// Price events
	case "price.created":
		return h.handlePriceCreated(event)
	case "price.deleted":
		return h.handlePriceDeleted(event)
	case "price.updated":
		return h.handlePriceUpdated(event)

	// Product events
	case "product.created":
		return h.handleProductCreated(event)
	case "product.deleted":
		return h.handleProductDeleted(event)
	case "product.updated":
		return h.handleProductUpdated(event)

	// Subscription schedule events
	case "subscription_schedule.aborted":
		return h.handleSubscriptionScheduleAborted(event)
	case "subscription_schedule.canceled":
		return h.handleSubscriptionScheduleCanceled(event)
	case "subscription_schedule.completed":
		return h.handleSubscriptionScheduleCompleted(event)
	case "subscription_schedule.created":
		return h.handleSubscriptionScheduleCreated(event)
	case "subscription_schedule.expiring":
		return h.handleSubscriptionScheduleExpiring(event)
	case "subscription_schedule.released":
		return h.handleSubscriptionScheduleReleased(event)
	case "subscription_schedule.updated":
		return h.handleSubscriptionScheduleUpdated(event)

	// Other potentially important events we'll support later
	case "customer.created":
		return h.handleCustomerCreated(event)
	case "customer.updated":
		return h.handleCustomerUpdated(event)
	case "customer.deleted":
		return h.handleCustomerDeleted(event)
	case "subscription.created":
		return h.handleSubscriptionCreated(event)
	case "subscription.updated":
		return h.handleSubscriptionUpdated(event)
	case "subscription.deleted":
		return h.handleSubscriptionDeleted(event)
	case "invoice.created":
		return h.handleInvoiceCreated(event)
	case "invoice.paid":
		return h.handleInvoicePaid(event)
	case "invoice.payment_failed":
		return h.handleInvoicePaymentFailed(event)

	default:
		h.logger.Info().
			Str("event_type", string(event.Type)).
			Msg("Received unhandled Stripe event type")
		return nil
	}
}

// Event handlers - stub implementations for all required events

// Checkout session handlers
func (h *StripeWebhookHandler) handleCheckoutSessionAsyncPaymentFailed(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing checkout.session.async_payment_failed")
	return nil
}

func (h *StripeWebhookHandler) handleCheckoutSessionAsyncPaymentSucceeded(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing checkout.session.async_payment_succeeded")
	return nil
}

func (h *StripeWebhookHandler) handleCheckoutSessionCompleted(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing checkout.session.completed")
	return nil
}

func (h *StripeWebhookHandler) handleCheckoutSessionExpired(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing checkout.session.expired")
	return nil
}

// Person handlers
func (h *StripeWebhookHandler) handlePersonCreated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing person.created")
	return nil
}

func (h *StripeWebhookHandler) handlePersonDeleted(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing person.deleted")
	return nil
}

func (h *StripeWebhookHandler) handlePersonUpdated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing person.updated")
	return nil
}

// Price handlers
func (h *StripeWebhookHandler) handlePriceCreated(event stripe.Event) error {
	// Parse the webhook payload
	var stripePrice stripe.Price
	err := json.Unmarshal(event.Data.Raw, &stripePrice)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal Stripe price data")
		return err
	}

	h.logger.Info().
		Str("stripe_price_id", stripePrice.ID).
		Str("stripe_product_id", stripePrice.Product.ID).
		Msg("Processing Stripe price.created event")

	// Check if this price already exists in our database
	ctx := context.Background()
	existingPrice, err := h.priceRepo.GetByStripeID(ctx, stripePrice.ID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_price_id", stripePrice.ID).
			Msg("Error checking for existing price")
		return err
	}

	// If price already exists, we don't need to create it
	if existingPrice != nil {
		h.logger.Info().
			Str("stripe_price_id", stripePrice.ID).
			Str("price_id", existingPrice.ID.String()).
			Msg("Price already exists in database, skipping creation")
		return nil
	}

	// Find the associated product
	product, err := h.productRepo.GetByStripeID(ctx, stripePrice.Product.ID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_product_id", stripePrice.Product.ID).
			Msg("Error finding associated product")
		return err
	}

	// If product doesn't exist, we need to fetch it from Stripe and create it
	if product == nil {
		h.logger.Warn().
			Str("stripe_product_id", stripePrice.Product.ID).
			Msg("Associated product not found, fetching from Stripe")

		// This would require implementing a method to fetch the product from Stripe
		// and create it in our database
		err = h.fetchAndCreateProduct(stripePrice.Product.ID)
		if err != nil {
			h.logger.Error().Err(err).
				Str("stripe_product_id", stripePrice.Product.ID).
				Msg("Failed to fetch and create product")
			return err
		}

		// Now try to get the product again
		product, err = h.productRepo.GetByStripeID(ctx, stripePrice.Product.ID)
		if err != nil || product == nil {
			h.logger.Error().Err(err).
				Str("stripe_product_id", stripePrice.Product.ID).
				Msg("Still unable to find product after creation attempt")
			return fmt.Errorf("unable to find or create associated product")
		}
	}

	// Create a new price model
	priceType := "one_time"
	var interval string
	var intervalCount int

	if stripePrice.Recurring != nil {
		priceType = "recurring"
		interval = string(stripePrice.Recurring.Interval)
		intervalCount = int(stripePrice.Recurring.IntervalCount)
	}

	newPrice := &model.Price{
		ID:            uuid.New(),
		ProductID:     product.ID,
		Name:          getPriceName(stripePrice, product.Name),
		Amount:        stripePrice.UnitAmount,
		Currency:      string(stripePrice.Currency),
		Type:          priceType,
		Interval:      interval,
		IntervalCount: intervalCount,
		Active:        stripePrice.Active,
		StripeID:      stripePrice.ID,
		CreatedAt:     time.Unix(stripePrice.Created, 0),
		UpdatedAt:     time.Now(),
	}

	// Save the price to our database
	err = h.priceRepo.Create(ctx, newPrice)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_price_id", stripePrice.ID).
			Msg("Failed to save price to database")
		return err
	}

	// Now we need to create a variant for this price
	// This is somewhat complex as we need to determine the option values
	// based on the price's metadata or product metadata
	variant, err := h.createVariantForPrice(ctx, newPrice, product, stripePrice)
	if err != nil {
		h.logger.Error().Err(err).
			Str("price_id", newPrice.ID.String()).
			Msg("Failed to create variant for price")
		// Continue anyway since the price was created
	}

	// Log success
	h.logger.Info().
		Str("stripe_price_id", stripePrice.ID).
		Str("price_id", newPrice.ID.String()).
		Str("product_id", product.ID.String()).
		Interface("variant", variant != nil).
		Msg("Successfully created price from Stripe webhook")

	return nil
}

// Helper to create a variant for a price
func (h *StripeWebhookHandler) createVariantForPrice(ctx context.Context, price *model.Price, product *model.Product, stripePrice stripe.Price) (*model.Variant, error) {
	// Extract options from metadata
	options := make(map[string]string)

	// First check price metadata
	if stripePrice.Metadata != nil {
		for key, value := range stripePrice.Metadata {
			if key == "weight" || key == "grind" || strings.HasPrefix(key, "option_") {
				optionKey := key
				if strings.HasPrefix(key, "option_") {
					optionKey = strings.TrimPrefix(key, "option_")
				}
				options[optionKey] = value
			}
		}
	}

	// If no options found, this might be a default variant
	isDefaultVariant := len(options) == 0

	// For default variants, check if one already exists
	if isDefaultVariant {
		variants, err := h.variantRepo.GetByProductID(ctx, product.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check for existing variants: %w", err)
		}

		// If there's already a default variant, we'll update its price
		for _, v := range variants {
			if len(v.Options) == 0 {
				// This is a default variant, update its price
				v.PriceID = price.ID
				v.StripePriceID = price.StripeID
				v.UpdatedAt = time.Now()

				err := h.variantRepo.Update(ctx, v)
				if err != nil {
					return nil, fmt.Errorf("failed to update default variant: %w", err)
				}

				h.logger.Info().
					Str("variant_id", v.ID.String()).
					Str("price_id", price.ID.String()).
					Msg("Updated existing default variant with new price")

				return v, nil
			}
		}
	}

	// Create a new variant
	variant := &model.Variant{
		ID:            uuid.New(),
		ProductID:     product.ID,
		PriceID:       price.ID,
		StripePriceID: price.StripeID,
		Weight:        product.Weight, // Default to product weight
		Options:       options,
		Active:        true,
		StockLevel:    product.StockLevel, // Default to product stock level
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Try to parse weight if present in options
	if weightStr, ok := options["weight"]; ok {
		weight := convertWeightToGrams(weightStr)
		if weight > 0 {
			variant.Weight = weight
		}
	}

	// Save the variant
	err := h.variantRepo.Create(ctx, variant)
	if err != nil {
		return nil, fmt.Errorf("failed to create variant: %w", err)
	}

	// Publish variant created event
	variantCreatedPayload := events.VariantCreatedPayload{
		VariantID:       variant.ID.String(),
		ProductID:       product.ID.String(),
		PriceID:         price.ID.String(),
		StripeProductID: product.StripeID,
		StripePriceID:   price.StripeID,
		Weight:          getOptionValue(options, "weight", fmt.Sprintf("%dg", product.Weight)),
		Grind:           getOptionValue(options, "grind", "Whole Bean"),
		OptionValues:    options,
		Amount:          price.Amount,
		Currency:        price.Currency,
		Active:          variant.Active,
		StockLevel:      variant.StockLevel,
		CreatedAt:       variant.CreatedAt,
	}

	err = h.eventBus.Publish(events.TopicVariantCreated, variantCreatedPayload)
	if err != nil {
		h.logger.Error().Err(err).
			Str("variant_id", variant.ID.String()).
			Msg("Failed to publish variant created event")
		// Continue anyway since the variant was created
	}

	return variant, nil
}

// Helper to generate a nice price name
func getPriceName(stripePrice stripe.Price, productName string) string {
	// If price has a nickname, use that
	if stripePrice.Nickname != "" {
		return stripePrice.Nickname
	}

	// Otherwise construct a name based on the product and price type
	name := productName

	// Add price details
	amount := float64(stripePrice.UnitAmount) / 100.0 // Convert cents to dollars
	currency := strings.ToUpper(string(stripePrice.Currency))

	if stripePrice.Recurring != nil {
		interval := string(stripePrice.Recurring.Interval)
		intervalCount := stripePrice.Recurring.IntervalCount

		// Format interval for readability
		intervalStr := interval
		if intervalCount > 1 {
			intervalStr = fmt.Sprintf("%d %ss", intervalCount, interval)
		} else if interval == "month" {
			intervalStr = "Monthly"
		} else if interval == "week" {
			intervalStr = "Weekly"
		} else if interval == "year" {
			intervalStr = "Annual"
		}

		return fmt.Sprintf("%s - %s (%.2f %s / %s)", name, intervalStr, amount, currency, interval)
	}

	return fmt.Sprintf("%s - One-time (%.2f %s)", name, amount, currency)
}

func (h *StripeWebhookHandler) handlePriceDeleted(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing price.deleted")
	return nil
}

func (h *StripeWebhookHandler) handlePriceUpdated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing price.updated")
	return nil
}

// Product handlers
// handleProductCreated processes a product.created webhook event
func (h *StripeWebhookHandler) handleProductCreated(event stripe.Event) error {
	// Parse the webhook payload
	var stripeProduct stripe.Product
	err := json.Unmarshal(event.Data.Raw, &stripeProduct)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal Stripe product data")
		return err
	}

	h.logger.Info().
		Str("stripe_product_id", stripeProduct.ID).
		Str("name", stripeProduct.Name).
		Msg("Processing Stripe product.created event")

	// Check if this product already exists in our database
	ctx := context.Background()

	// Use a repository to find the product by Stripe ID
	// This would require access to your repositories
	existingProduct, err := h.productRepo.GetByStripeID(ctx, stripeProduct.ID)
	if err != nil {
		h.logger.Error().Err(err).Str("stripe_product_id", stripeProduct.ID).
			Msg("Error checking for existing product")
		return err
	}

	// If product already exists, we don't need to create it
	if existingProduct != nil {
		h.logger.Info().
			Str("stripe_product_id", stripeProduct.ID).
			Str("product_id", existingProduct.ID.String()).
			Msg("Product already exists in database, skipping creation")
		return nil
	}

	// Extract metadata
	options := make(map[string][]string)
	allowSubscription := false
	var origin, roastLevel, flavorNotes string
	weight := 340 // Default weight in grams

	// Parse metadata
	if stripeProduct.Metadata != nil {
		// Check for origin
		if val, ok := stripeProduct.Metadata["origin"]; ok {
			origin = val
		}

		// Check for roast level
		if val, ok := stripeProduct.Metadata["roast_level"]; ok {
			roastLevel = val
		}

		// Check for flavor notes
		if val, ok := stripeProduct.Metadata["flavor_notes"]; ok {
			flavorNotes = val
		}

		// Check for weight
		if val, ok := stripeProduct.Metadata["weight"]; ok {
			if w, err := strconv.Atoi(val); err == nil {
				weight = w
			}
		}

		// Check for subscription flag
		if val, ok := stripeProduct.Metadata["allow_subscription"]; ok {
			allowSubscription = val == "true"
		}

		// Check for options
		for key, val := range stripeProduct.Metadata {
			// Options are stored with a prefix to distinguish them
			if strings.HasPrefix(key, "option_") {
				optionKey := strings.TrimPrefix(key, "option_")
				optionValues := strings.Split(val, ",")
				options[optionKey] = optionValues
			}
		}
	}

	// Create a new product model
	newProduct := &model.Product{
		ID:                uuid.New(),
		StripeID:          stripeProduct.ID,
		Name:              stripeProduct.Name,
		Description:       stripeProduct.Description,
		ImageURL:          getFirstImage(stripeProduct.Images),
		Origin:            origin,
		RoastLevel:        roastLevel,
		FlavorNotes:       flavorNotes,
		Active:            stripeProduct.Active,
		Archived:          false,
		AllowSubscription: allowSubscription,
		StockLevel:        0, // Default to 0 until we know the actual stock
		Weight:            weight,
		Options:           options,
		CreatedAt:         time.Unix(stripeProduct.Created, 0),
		UpdatedAt:         time.Now(),
	}

	// Save the product to our database
	err = h.productRepo.Create(ctx, newProduct)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_product_id", stripeProduct.ID).
			Msg("Failed to save product to database")
		return err
	}

	// Publish event that a product was created
	payload := events.ProductCreatedPayload{
		ProductID:         newProduct.ID.String(),
		Name:              newProduct.Name,
		Description:       newProduct.Description,
		ImageURL:          newProduct.ImageURL,
		StockLevel:        newProduct.StockLevel,
		Weight:            newProduct.Weight,
		Origin:            newProduct.Origin,
		RoastLevel:        newProduct.RoastLevel,
		FlavorNotes:       newProduct.FlavorNotes,
		Options:           newProduct.Options,
		AllowSubscription: newProduct.AllowSubscription,
		Active:            newProduct.Active,
		CreatedAt:         newProduct.CreatedAt,
	}

	err = h.eventBus.Publish(events.TopicProductCreated, payload)
	if err != nil {
		h.logger.Error().Err(err).
			Str("product_id", newProduct.ID.String()).
			Msg("Failed to publish product created event")
		// Don't return error since product is already saved
	}

	h.logger.Info().
		Str("stripe_product_id", stripeProduct.ID).
		Str("product_id", newProduct.ID.String()).
		Msg("Successfully created product from Stripe webhook")

	return nil
}

// Helper to get first image from a slice
func getFirstImage(images []string) string {
	if len(images) > 0 {
		return images[0]
	}
	return ""
}

// handleProductDeleted processes a product.deleted webhook event
func (h *StripeWebhookHandler) handleProductDeleted(event stripe.Event) error {
	// Parse the webhook payload
	var stripeProduct stripe.Product
	err := json.Unmarshal(event.Data.Raw, &stripeProduct)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal Stripe product data")
		return err
	}

	h.logger.Info().
		Str("stripe_product_id", stripeProduct.ID).
		Msg("Processing Stripe product.deleted event")

	// Find the product in our database
	ctx := context.Background()
	existingProduct, err := h.productRepo.GetByStripeID(ctx, stripeProduct.ID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_product_id", stripeProduct.ID).
			Msg("Error checking for existing product")
		return err
	}

	// If product doesn't exist, nothing to do
	if existingProduct == nil {
		h.logger.Info().
			Str("stripe_product_id", stripeProduct.ID).
			Msg("Product not found in database, nothing to delete")
		return nil
	}

	// For safety, we don't actually delete the product, we archive it
	err = h.productRepo.Archive(ctx, existingProduct.ID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_product_id", stripeProduct.ID).
			Str("product_id", existingProduct.ID.String()).
			Msg("Failed to archive product in database")
		return err
	}

	// Publish product deleted event
	payload := map[string]interface{}{
		"product_id":    existingProduct.ID.String(),
		"stripe_id":     existingProduct.StripeID,
		"deleted_at":    time.Now(),
		"delete_source": "stripe_webhook",
	}

	err = h.eventBus.Publish(events.TopicProductDeleted, payload)
	if err != nil {
		h.logger.Error().Err(err).
			Str("product_id", existingProduct.ID.String()).
			Msg("Failed to publish product deleted event")
		// Don't return error since product is already archived
	}

	h.logger.Info().
		Str("stripe_product_id", stripeProduct.ID).
		Str("product_id", existingProduct.ID.String()).
		Msg("Successfully archived product from Stripe webhook")

	return nil
}

// handleProductUpdated processes a product.updated webhook event
func (h *StripeWebhookHandler) handleProductUpdated(event stripe.Event) error {
	// Parse the webhook payload
	var stripeProduct stripe.Product
	err := json.Unmarshal(event.Data.Raw, &stripeProduct)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal Stripe product data")
		return err
	}

	h.logger.Info().
		Str("stripe_product_id", stripeProduct.ID).
		Str("name", stripeProduct.Name).
		Msg("Processing Stripe product.updated event")

	// Find the product in our database
	ctx := context.Background()
	existingProduct, err := h.productRepo.GetByStripeID(ctx, stripeProduct.ID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_product_id", stripeProduct.ID).
			Msg("Error checking for existing product")
		return err
	}

	// If product doesn't exist, create it
	if existingProduct == nil {
		h.logger.Warn().
			Str("stripe_product_id", stripeProduct.ID).
			Msg("Product not found in database, handling as creation instead")
		return h.handleProductCreated(event)
	}

	// Update the product with new values
	existingProduct.Name = stripeProduct.Name
	existingProduct.Description = stripeProduct.Description
	existingProduct.Active = stripeProduct.Active
	existingProduct.UpdatedAt = time.Now()

	// Update image if available
	if len(stripeProduct.Images) > 0 {
		existingProduct.ImageURL = stripeProduct.Images[0]
	}

	// Update metadata fields if present
	if stripeProduct.Metadata != nil {
		if val, ok := stripeProduct.Metadata["origin"]; ok {
			existingProduct.Origin = val
		}
		if val, ok := stripeProduct.Metadata["roast_level"]; ok {
			existingProduct.RoastLevel = val
		}
		if val, ok := stripeProduct.Metadata["flavor_notes"]; ok {
			existingProduct.FlavorNotes = val
		}
		if val, ok := stripeProduct.Metadata["weight"]; ok {
			if w, err := strconv.Atoi(val); err == nil {
				existingProduct.Weight = w
			}
		}
		if val, ok := stripeProduct.Metadata["allow_subscription"]; ok {
			existingProduct.AllowSubscription = val == "true"
		}

		// Update options if present
		// This is a bit more complex as we need to handle option removals
		newOptions := make(map[string][]string)
		for key, val := range stripeProduct.Metadata {
			if strings.HasPrefix(key, "option_") {
				optionKey := strings.TrimPrefix(key, "option_")
				optionValues := strings.Split(val, ",")
				newOptions[optionKey] = optionValues
			}
		}

		// Only update options if we found some in the metadata
		if len(newOptions) > 0 {
			existingProduct.Options = newOptions
		}
	}

	// Save the updated product
	err = h.productRepo.Update(ctx, existingProduct)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_product_id", stripeProduct.ID).
			Str("product_id", existingProduct.ID.String()).
			Msg("Failed to update product in database")
		return err
	}

	// Publish product updated event
	payload := map[string]interface{}{
		"product_id":    existingProduct.ID.String(),
		"stripe_id":     existingProduct.StripeID,
		"name":          existingProduct.Name,
		"updated_at":    existingProduct.UpdatedAt,
		"update_source": "stripe_webhook",
	}

	err = h.eventBus.Publish(events.TopicProductUpdated, payload)
	if err != nil {
		h.logger.Error().Err(err).
			Str("product_id", existingProduct.ID.String()).
			Msg("Failed to publish product updated event")
		// Don't return error since product is already updated
	}

	h.logger.Info().
		Str("stripe_product_id", stripeProduct.ID).
		Str("product_id", existingProduct.ID.String()).
		Msg("Successfully updated product from Stripe webhook")

	return nil
}

// Subscription schedule handlers
func (h *StripeWebhookHandler) handleSubscriptionScheduleAborted(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing subscription_schedule.aborted")
	return nil
}

func (h *StripeWebhookHandler) handleSubscriptionScheduleCanceled(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing subscription_schedule.canceled")
	return nil
}

func (h *StripeWebhookHandler) handleSubscriptionScheduleCompleted(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing subscription_schedule.completed")
	return nil
}

func (h *StripeWebhookHandler) handleSubscriptionScheduleCreated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing subscription_schedule.created")
	return nil
}

func (h *StripeWebhookHandler) handleSubscriptionScheduleExpiring(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing subscription_schedule.expiring")
	return nil
}

func (h *StripeWebhookHandler) handleSubscriptionScheduleReleased(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing subscription_schedule.released")
	return nil
}

func (h *StripeWebhookHandler) handleSubscriptionScheduleUpdated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing subscription_schedule.updated")
	return nil
}

// Additional event handlers for common Stripe events
func (h *StripeWebhookHandler) handleCustomerCreated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing customer.created")
	return nil
}

func (h *StripeWebhookHandler) handleCustomerUpdated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing customer.updated")
	return nil
}

func (h *StripeWebhookHandler) handleCustomerDeleted(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing customer.deleted")
	return nil
}

func (h *StripeWebhookHandler) handleSubscriptionCreated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing subscription.created")
	return nil
}

func (h *StripeWebhookHandler) handleSubscriptionUpdated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing subscription.updated")
	return nil
}

func (h *StripeWebhookHandler) handleSubscriptionDeleted(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing subscription.deleted")
	return nil
}

func (h *StripeWebhookHandler) handleInvoiceCreated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing invoice.created")
	return nil
}

func (h *StripeWebhookHandler) handleInvoicePaid(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing invoice.paid")
	return nil
}

func (h *StripeWebhookHandler) handleInvoicePaymentFailed(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing invoice.payment_failed")
	return nil
}

// Helper function to handle fetching and creating a product from Stripe
func (h *StripeWebhookHandler) fetchAndCreateProduct(stripeProductID string) error {
	// In a real implementation, you would:
	// 1. Use the Stripe API to fetch the product details
	// 2. Create a product record in your database
	// For now, return an error as this is just a stub
	return fmt.Errorf("fetchAndCreateProduct not implemented")
}

// Helper function to convert weight strings to grams
func convertWeightToGrams(weightStr string) int {
	// Common conversion factors
	const (
		gramsPerOz = 28
		gramsPerLb = 454
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

	// Try to parse direct gram values
	if strings.HasSuffix(weightStr, "g") {
		value := strings.TrimSuffix(weightStr, "g")
		if g, err := strconv.Atoi(value); err == nil {
			return g
		}
	}

	// Default to 1 gram if we can't parse
	return 1
}

// Helper to safely get option values with a default fallback
func getOptionValue(options map[string]string, key, defaultValue string) string {
	if value, exists := options[key]; exists {
		return value
	}
	return defaultValue
}
