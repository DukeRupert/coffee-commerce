// internal/api/handler/admin_handler.go
package handler

import (
	"net/http"

	"github.com/dukerupert/coffee-commerce/internal/api"
	"github.com/dukerupert/coffee-commerce/internal/interfaces"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AdminHandler interface {
	SyncStripeProductIDs(c echo.Context) error
	HealthCheck(c echo.Context) error
}

// adminHandler handles administrative operations
type adminHandler struct {
	logger       zerolog.Logger
	priceService interfaces.PriceService
	productRepo  interfaces.ProductRepository
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(
	logger *zerolog.Logger,
	priceService interfaces.PriceService,
	productRepo interfaces.ProductRepository,
) *adminHandler {
	sublogger := logger.With().Str("component", "admin_handler").Logger()
	return &adminHandler{
		logger:       sublogger,
		priceService: priceService,
		productRepo:  productRepo,
	}
}

// SyncStripeProductIDs handles POST /api/admin/sync-stripe-ids
func (h *adminHandler) SyncStripeProductIDs(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	h.logger.Info().
		Str("handler", "AdminHandler.SyncStripeProductIDs").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Msg("Handling Stripe product ID sync request")

	// Call the sync method on the price service
	syncResults, err := h.priceService.SyncStripeProductIDs(ctx)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Msg("Failed to sync Stripe product IDs")

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to sync Stripe product IDs",
			Code:    "SYNC_FAILED",
		})
	}

	h.logger.Info().
		Str("request_id", requestID).
		Interface("sync_results", syncResults).
		Msg("Stripe product ID sync completed successfully")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Stripe product ID sync completed successfully",
		"results": syncResults,
	})
}

// HealthCheck handles GET /api/admin/health
func (h *adminHandler) HealthCheck(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	h.logger.Debug().
		Str("handler", "AdminHandler.HealthCheck").
		Str("request_id", requestID).
		Msg("Handling health check request")

	// Get product count for basic health info
	_, total, err := h.productRepo.List(ctx, 0, 1, true, true)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to check database health")
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"status": "unhealthy",
			"error":  "Database connection failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":        "healthy",
		"total_products": total,
		"database":      "connected",
		"timestamp":     ctx.Value("timestamp"),
	})
}