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
	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/dukerupert/coffee-commerce/internal/events"
	"github.com/dukerupert/coffee-commerce/internal/interfaces"
	"github.com/dukerupert/coffee-commerce/internal/sync"
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
	syncRepo     interfaces.SyncHashRepository
}

func NewStripeWebhookHandler(
	logger *zerolog.Logger,
	stripeConfig *config.StripeConfig,
	eventBus events.EventBus, productRepo interfaces.ProductRepository, priceRepo interfaces.PriceRepository,
	variantRepo interfaces.VariantRepository, syncRepo interfaces.SyncHashRepository) *StripeWebhookHandler {

	return &StripeWebhookHandler{
		logger:       logger.With().Str("component", "stripe_webhook_handler").Logger(),
		stripeConfig: stripeConfig,
		eventBus:     eventBus,
		productRepo:  productRepo,
		priceRepo:    priceRepo,
		variantRepo:  variantRepo,
		syncRepo:     syncRepo,
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
		return c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Status:  http.StatusServiceUnavailable,
			Message: "Failed to read request body",
		})
	}

	// Get the signature header
	signatureHeader := c.Request().Header.Get("Stripe-Signature")
	if signatureHeader == "" {
		h.logger.Warn().Msg("Missing Stripe-Signature header")
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Missing Stripe-Signature header",
		})
	}

	// Verify the signature
	event, err := webhook.ConstructEvent(body, signatureHeader, h.stripeConfig.WebhookSecret)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to verify webhook signature")
		return c.JSON(http.StatusBadRequest, ErrorResponse{
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

	// Find the variant that corresponds to this Stripe product
	existingVariant, err := h.variantRepo.GetByStripeProductID(ctx, stripePrice.Product.ID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_product_id", stripePrice.Product.ID).
			Msg("Error finding associated variant")
		return err
	}

	// If variant doesn't exist, we need to handle this case
	if existingVariant == nil {
		h.logger.Warn().
			Str("stripe_product_id", stripePrice.Product.ID).
			Str("stripe_price_id", stripePrice.ID).
			Msg("No variant found for Stripe product - this price belongs to a Stripe product not in our system")

		// We could try to fetch the Stripe product and create a variant, but for now we'll skip
		h.logger.Info().
			Str("stripe_product_id", stripePrice.Product.ID).
			Str("stripe_price_id", stripePrice.ID).
			Msg("Skipping price creation - no associated variant found")
		return nil
	}

	// Get the parent product for the price record
	parentProduct, err := h.productRepo.GetByID(ctx, existingVariant.ProductID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("product_id", existingVariant.ProductID.String()).
			Msg("Error finding parent product for variant")
		return err
	}

	if parentProduct == nil {
		h.logger.Error().
			Str("product_id", existingVariant.ProductID.String()).
			Str("variant_id", existingVariant.ID.String()).
			Msg("Parent product not found for variant")
		return fmt.Errorf("parent product not found for variant %s", existingVariant.ID.String())
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
		ProductID:     parentProduct.ID, // Reference the parent product
		Name:          getPriceName(stripePrice, parentProduct.Name),
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

	// Update the variant to reference this new price
	existingVariant.PriceID = newPrice.ID
	existingVariant.StripePriceID = stripePrice.ID
	existingVariant.UpdatedAt = time.Now()

	err = h.variantRepo.Update(ctx, existingVariant)
	if err != nil {
		h.logger.Error().Err(err).
			Str("variant_id", existingVariant.ID.String()).
			Str("price_id", newPrice.ID.String()).
			Msg("Failed to update variant with new price")
		// Continue anyway since the price was created
	} else {
		h.logger.Info().
			Str("variant_id", existingVariant.ID.String()).
			Str("price_id", newPrice.ID.String()).
			Msg("Updated variant with new price")
	}

	// Publish variant updated event (since the variant now has a new price)
	variantUpdatedPayload := events.VariantUpdatedPayload{
		VariantID:       existingVariant.ID.String(),
		ProductID:       parentProduct.ID.String(),
		PriceID:         newPrice.ID.String(),
		StripeProductID: existingVariant.StripeProductID,
		StripePriceID:   newPrice.StripeID,
		Amount:          newPrice.Amount,
		Currency:        newPrice.Currency,
		PriceType:       newPrice.Type,
		Interval:        newPrice.Interval,
		IntervalCount:   newPrice.IntervalCount,
		UpdatedAt:       time.Now(),
		UpdateSource:    "stripe_webhook",
	}

	err = h.eventBus.Publish(events.TopicVariantUpdated, variantUpdatedPayload)
	if err != nil {
		h.logger.Error().Err(err).
			Str("variant_id", existingVariant.ID.String()).
			Msg("Failed to publish variant updated event")
		// Don't return error since the price and variant were updated successfully
	}

	// Log success
	h.logger.Info().
		Str("stripe_price_id", stripePrice.ID).
		Str("price_id", newPrice.ID.String()).
		Str("variant_id", existingVariant.ID.String()).
		Str("product_id", parentProduct.ID.String()).
		Msg("Successfully created price and updated variant from Stripe webhook")

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

	// Check if this variant already exists in our database
	ctx := context.Background()

	// Find variant with this Stripe Product ID
	existingVariant, err := h.variantRepo.GetByStripeID(ctx, stripeProduct.ID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_product_id", stripeProduct.ID).
			Msg("Error checking for existing variant")
		return err
	}

	// If variant already exists, we don't need to create it
	if existingVariant != nil {
		h.logger.Info().
			Str("stripe_product_id", stripeProduct.ID).
			Str("variant_id", existingVariant.ID.String()).
			Msg("Variant already exists in database, skipping creation")
		return nil
	}

	// Extract metadata to determine which product this variant belongs to
	var productID uuid.UUID
	weight := 340 // Default weight in grams
	options := make(map[string]string)

	// Parse metadata for product association and variant options
	if stripeProduct.Metadata != nil {
		// Check for product ID - this is crucial to link variant to product
		if val, ok := stripeProduct.Metadata["product_id"]; ok {
			if pid, err := uuid.Parse(val); err == nil {
				productID = pid
			} else {
				h.logger.Error().Err(err).
					Str("stripe_product_id", stripeProduct.ID).
					Str("product_id_value", val).
					Msg("Invalid product_id in metadata")
				return fmt.Errorf("invalid product_id in metadata: %w", err)
			}
		} else {
			// Try to find product ID from original_product_id
			if val, ok := stripeProduct.Metadata["original_product_id"]; ok {
				if pid, err := uuid.Parse(val); err == nil {
					productID = pid
				}
			}
		}

		// If we still don't have a product ID, we can't create the variant
		if productID == uuid.Nil {
			h.logger.Error().
				Str("stripe_product_id", stripeProduct.ID).
				Msg("No product_id found in metadata, cannot create variant")
			return fmt.Errorf("no product_id in metadata for Stripe product %s", stripeProduct.ID)
		}

		// Check for weight
		if val, ok := stripeProduct.Metadata["weight"]; ok {
			if w, err := strconv.Atoi(val); err == nil {
				weight = w
			}
		}

		// Extract all option values from metadata
		for key, val := range stripeProduct.Metadata {
			// Skip special keys used for product association
			if key == "product_id" || key == "original_product_id" {
				continue
			}

			// Add to options map
			options[key] = val
		}
	}

	// Verify that the product exists
	product, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("product_id", productID.String()).
			Msg("Error retrieving product")
		return err
	}

	if product == nil {
		h.logger.Error().
			Str("product_id", productID.String()).
			Msg("Product not found, cannot create variant")
		return fmt.Errorf("product with ID %s not found", productID)
	}

	// Create a price for this variant (default price)
	// Note: Normally this would be created from a Stripe price webhook
	// but we'll create a temporary one
	priceID := uuid.New()
	price := &model.Price{
		ID:        priceID,
		ProductID: productID,
		Name:      stripeProduct.Name + " - Default Price",
		Amount:    1000, // $10.00 default
		Currency:  "USD",
		Type:      "one_time",
		Active:    true,
		StripeID:  "temp_" + uuid.New().String(), // This would normally be a Stripe price ID
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the price
	err = h.priceRepo.Create(ctx, price)
	if err != nil {
		h.logger.Error().Err(err).
			Str("product_id", productID.String()).
			Msg("Failed to create price for variant")
		return err
	}

	// Create the new variant
	newVariant := &model.Variant{
		ID:              uuid.New(),
		ProductID:       productID,
		PriceID:         priceID,
		StripeProductID: stripeProduct.ID,
		StripePriceID:   price.StripeID, // This would normally be set from a price webhook
		Weight:          weight,
		Options:         options,
		Active:          stripeProduct.Active,
		StockLevel:      0, // Default to 0 until we know the actual stock
		CreatedAt:       time.Unix(stripeProduct.Created, 0),
		UpdatedAt:       time.Now(),
	}

	// Save the variant to our database
	err = h.variantRepo.Create(ctx, newVariant)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_product_id", stripeProduct.ID).
			Str("product_id", productID.String()).
			Msg("Failed to save variant to database")
		return err
	}

	// Publish event that a variant was created
	payload := events.VariantCreatedPayload{
		VariantID:       newVariant.ID.String(),
		ProductID:       productID.String(),
		PriceID:         priceID.String(),
		StripeProductID: stripeProduct.ID,
		StripePriceID:   price.StripeID,
		Weight:          fmt.Sprintf("%dg", weight),
		OptionValues:    options,
		Amount:          price.Amount,
		Currency:        price.Currency,
		Active:          newVariant.Active,
		StockLevel:      newVariant.StockLevel,
		CreatedAt:       newVariant.CreatedAt,
	}

	err = h.eventBus.Publish(events.TopicVariantCreated, payload)
	if err != nil {
		h.logger.Error().Err(err).
			Str("variant_id", newVariant.ID.String()).
			Msg("Failed to publish variant created event")
		// Don't return error since variant is already saved
	}

	h.logger.Info().
		Str("stripe_product_id", stripeProduct.ID).
		Str("variant_id", newVariant.ID.String()).
		Str("product_id", productID.String()).
		Msg("Successfully created variant from Stripe product webhook")

	return nil
}

// handleProductDeleted processes a product.deleted webhook event
// Note: In our system, Stripe "products" map to variants, not products
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

	// Find the variant in our database by Stripe product ID
	ctx := context.Background()

	// First, we need to find all variants and filter by StripeProductID
	// since we don't have a direct GetByStripeProductID method
	existingVariant, err := h.variantRepo.GetByStripeProductID(ctx, stripeProduct.ID)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_product_id", stripeProduct.ID).
			Msg("Error checking for existing variant")
		return err
	}

	// If variant doesn't exist, nothing to do
	if existingVariant == nil {
		h.logger.Info().
			Str("stripe_product_id", stripeProduct.ID).
			Msg("Variant not found in database, nothing to delete")
		return nil
	}

	// Get the parent product for logging and event publishing
	parentProduct, err := h.productRepo.GetByID(ctx, existingVariant.ProductID)
	if err != nil {
		h.logger.Warn().Err(err).
			Str("product_id", existingVariant.ProductID.String()).
			Msg("Could not retrieve parent product, continuing with variant deletion")
		// Continue anyway - we can still delete the variant
	}

	// For safety, we don't actually delete the variant, we deactivate it
	existingVariant.Active = false
	existingVariant.UpdatedAt = time.Now()

	err = h.variantRepo.Update(ctx, existingVariant)
	if err != nil {
		h.logger.Error().Err(err).
			Str("stripe_product_id", stripeProduct.ID).
			Str("variant_id", existingVariant.ID.String()).
			Msg("Failed to deactivate variant in database")
		return err
	}

	// Publish variant deleted event
	payload := map[string]interface{}{
		"variant_id":        existingVariant.ID.String(),
		"product_id":        existingVariant.ProductID.String(),
		"stripe_product_id": stripeProduct.ID,
		"stripe_price_id":   existingVariant.StripePriceID,
		"deleted_at":        time.Now(),
		"delete_source":     "stripe_webhook",
	}

	// Include parent product info if available
	if parentProduct != nil {
		payload["product_name"] = parentProduct.Name
	}

	err = h.eventBus.Publish(events.TopicVariantDeleted, payload)
	if err != nil {
		h.logger.Error().Err(err).
			Str("variant_id", existingVariant.ID.String()).
			Msg("Failed to publish variant deleted event")
		// Don't return error since variant is already deactivated
	}

	h.logger.Info().
		Str("stripe_product_id", stripeProduct.ID).
		Str("variant_id", existingVariant.ID.String()).
		Str("product_id", existingVariant.ProductID.String()).
		Msg("Successfully deactivated variant from Stripe webhook")

	return nil
}

// handleProductUpdated processes a product.updated webhook event
// Note: In our system, Stripe "products" map to variants, not products
func (h *StripeWebhookHandler) handleProductUpdated(event stripe.Event) error {
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

	ctx := context.Background()

	// Find existing variant
	existingVariant, err := h.variantRepo.GetByStripeProductID(ctx, stripeProduct.ID)
	if err != nil || existingVariant == nil {
		return h.handleProductCreated(event)
	}

	// Compute incoming hash
	incomingHash, err := sync.ComputeStripeProductHash(stripeProduct)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to compute incoming hash")
		return err
	}

	// Get stored hash
	storedSyncHash, err := h.syncRepo.GetByVariantAndStripeID(ctx, existingVariant.ID, stripeProduct.ID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get stored hash")
		return err
	}

	// Compare hashes
	if storedSyncHash != nil && storedSyncHash.ContentHash == incomingHash {
		h.logger.Debug().
			Str("stripe_product_id", stripeProduct.ID).
			Str("hash", incomingHash).
			Msg("Hash unchanged, skipping update")
		return nil
	}

	// Perform update
	err = h.updateVariantFromStripeProduct(ctx, existingVariant, stripeProduct)
	if err != nil {
		return err
	}

	// Store new hash
	newSyncHash := sync.CreateSyncHashRecord(
		existingVariant.ID,
		stripeProduct.ID,
		incomingHash,
		model.SyncSourceStripeWebhook,
	)

	err = h.syncRepo.Upsert(ctx, newSyncHash)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to store sync hash")
		// Don't fail the webhook for this
	}

	return nil
}

// updateVariantFromStripeProduct updates a variant with data from a Stripe product
func (h *StripeWebhookHandler) updateVariantFromStripeProduct(ctx context.Context, variant *model.Variant, stripeProduct stripe.Product) error {
	// Track what fields are being updated for logging
	updatedFields := []string{}

	// Update basic variant fields
	if variant.Active != stripeProduct.Active {
		variant.Active = stripeProduct.Active
		updatedFields = append(updatedFields, "active")
	}

	// Update variant-specific fields from Stripe product metadata
	if stripeProduct.Metadata != nil {
		// Update weight if present in metadata
		if val, ok := stripeProduct.Metadata["weight"]; ok {
			newWeight := convertWeightToGrams(val)
			if newWeight > 0 && newWeight != variant.Weight {
				variant.Weight = newWeight
				updatedFields = append(updatedFields, "weight")
			}
		}

		// Update stock level if present in metadata
		if val, ok := stripeProduct.Metadata["stock_level"]; ok {
			if stockLevel, err := strconv.Atoi(val); err == nil && stockLevel != variant.StockLevel {
				variant.StockLevel = stockLevel
				updatedFields = append(updatedFields, "stock_level")
			}
		}

		// Update variant options from metadata
		newOptions := make(map[string]string)
		optionsUpdated := false

		// Start with existing options
		for key, val := range variant.Options {
			newOptions[key] = val
		}

		// Look for option updates in metadata
		for key, val := range stripeProduct.Metadata {
			// Handle direct option fields
			if key == "weight" || key == "grind" {
				if existingValue, exists := newOptions[key]; !exists || existingValue != val {
					newOptions[key] = val
					optionsUpdated = true
				}
			}

			// Handle prefixed option fields (e.g., "variant_weight", "variant_grind")
			if strings.HasPrefix(key, "variant_") {
				optionKey := strings.TrimPrefix(key, "variant_")
				if existingValue, exists := newOptions[optionKey]; !exists || existingValue != val {
					newOptions[optionKey] = val
					optionsUpdated = true
				}
			}
		}

		// Update options if any changed
		if optionsUpdated {
			variant.Options = newOptions
			updatedFields = append(updatedFields, "options")
		}
	}

	// Update timestamp
	variant.UpdatedAt = time.Now()

	// Save the updated variant
	err := h.variantRepo.Update(ctx, variant)
	if err != nil {
		h.logger.Error().Err(err).
			Str("variant_id", variant.ID.String()).
			Str("stripe_product_id", stripeProduct.ID).
			Msg("Failed to update variant")
		return fmt.Errorf("failed to update variant: %w", err)
	}

	// Get parent product for event publishing
	parentProduct, err := h.productRepo.GetByID(ctx, variant.ProductID)
	if err != nil {
		h.logger.Warn().Err(err).
			Str("product_id", variant.ProductID.String()).
			Msg("Could not retrieve parent product for event publishing")
		// Continue anyway - variant update succeeded
	}

	// Publish variant updated event
	variantUpdatedPayload := events.VariantUpdatedPayload{
		VariantID:       variant.ID.String(),
		ProductID:       parentProduct.ID.String(),
		PriceID:         variant.PriceID.String(),
		StripeProductID: variant.StripeProductID,
		StripePriceID:   variant.StripePriceID,
		UpdatedAt:       variant.UpdatedAt,
		UpdateSource:    model.SyncSourceStripeWebhook,
	}

	// Get price information for the payload if available
	if variant.PriceID != uuid.Nil {
		price, err := h.priceRepo.GetByID(ctx, variant.PriceID)
		if err == nil && price != nil {
			variantUpdatedPayload.Amount = price.Amount
			variantUpdatedPayload.Currency = price.Currency
			variantUpdatedPayload.PriceType = price.Type
			variantUpdatedPayload.Interval = price.Interval
			variantUpdatedPayload.IntervalCount = price.IntervalCount
		}
	}

	err = h.eventBus.Publish(events.TopicVariantUpdated, variantUpdatedPayload)
	if err != nil {
		h.logger.Error().Err(err).
			Str("variant_id", variant.ID.String()).
			Msg("Failed to publish variant updated event")
		// Don't return error since variant update succeeded
	}

	// Log successful update
	h.logger.Info().
		Str("variant_id", variant.ID.String()).
		Str("stripe_product_id", stripeProduct.ID).
		Str("product_id", variant.ProductID.String()).
		Strs("updated_fields", updatedFields).
		Msg("Successfully updated variant from Stripe product")

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
