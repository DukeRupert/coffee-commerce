// internal/api/handler/stripe_webhook_handler.go
package handler

import (
	"io"
	"net/http"

	"github.com/dukerupert/coffee-commerce/config"
	"github.com/dukerupert/coffee-commerce/internal/api"
	events "github.com/dukerupert/coffee-commerce/internal/event"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripeWebhookHandler struct {
	logger       zerolog.Logger
	stripeConfig *config.StripeConfig
	eventBus     events.EventBus
}

func NewStripeWebhookHandler(
	logger *zerolog.Logger, 
	stripeConfig *config.StripeConfig,
	eventBus events.EventBus) *StripeWebhookHandler {
	
	return &StripeWebhookHandler{
		logger:       logger.With().Str("component", "stripe_webhook_handler").Logger(),
		stripeConfig: stripeConfig,
		eventBus:     eventBus,
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
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing price.created")
	
	// Here we would typically:
	// 1. Extract price data from event
	// 2. Check if this price is already tracked in our system
	// 3. If not, create a record in our database
	// 4. Publish an internal event that a price was created
	
	return nil
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
func (h *StripeWebhookHandler) handleProductCreated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing product.created")
	
	// Here we would typically:
	// 1. Extract product data from event
	// 2. Check if we're tracking this product in our system
	// 3. If not (e.g., created via Stripe dashboard), create a corresponding record
	// 4. Publish an internal event for other services
	
	return nil
}

func (h *StripeWebhookHandler) handleProductDeleted(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing product.deleted")
	return nil
}

func (h *StripeWebhookHandler) handleProductUpdated(event stripe.Event) error {
	h.logger.Debug().Interface("data", event.Data).Msg("Stub: Processing product.updated")
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