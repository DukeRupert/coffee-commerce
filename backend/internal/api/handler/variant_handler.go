// internal/api/handler/variant_handler.go
package handler

import (
	"net/http"

	"github.com/dukerupert/coffee-commerce/internal/api"
	interfaces "github.com/dukerupert/coffee-commerce/internal/repository/interface"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type VariantHandler interface {
	ListByProduct(c echo.Context) error
}

// variantHandler handles HTTP requests for variants
type variantHandler struct {
	logger      zerolog.Logger
	variantRepo interfaces.VariantRepository
	productRepo interfaces.ProductRepository
}

// NewVariantHandler creates a new variant handler
func NewVariantHandler(logger *zerolog.Logger, variantRepo interfaces.VariantRepository, productRepo interfaces.ProductRepository) *variantHandler {
	sublogger := logger.With().Str("component", "variant_handler").Logger()
	return &variantHandler{
		logger:      sublogger,
		variantRepo: variantRepo,
		productRepo: productRepo,
	}
}

// ListByProduct handles GET /api/products/:id/variants
func (h *variantHandler) ListByProduct(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// Parse product ID from URL
	productIDParam := c.Param("id")

	h.logger.Info().
		Str("handler", "VariantHandler.ListByProduct").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Str("product_id", productIDParam).
		Msg("Handling list variants by product ID request")

	// Convert string ID to UUID
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Str("product_id", productIDParam).
			Msg("Invalid product ID format")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid product ID format",
			Code:    "INVALID_ID_FORMAT",
		})
	}

	// Check if the product exists
	product, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("product_id", productID.String()).
			Msg("Error checking if product exists")

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve product information",
			Code:    "INTERNAL_ERROR",
		})
	}

	if product == nil {
		h.logger.Warn().
			Str("request_id", requestID).
			Str("product_id", productID.String()).
			Msg("Product not found")

		return c.JSON(http.StatusNotFound, api.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "Product not found",
			Code:    "PRODUCT_NOT_FOUND",
		})
	}

	// Get variants for the product
	variants, err := h.variantRepo.GetByProductID(ctx, productID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("product_id", productID.String()).
			Msg("Failed to retrieve variants")

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve variants",
			Code:    "INTERNAL_ERROR",
		})
	}

	// Convert to DTO if needed (we'll use the model directly for simplicity)
	// If the response contains sensitive information, you may want to map to a DTO

	// Return the variants
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": variants,
		"product": map[string]interface{}{
			"id":   product.ID,
			"name": product.Name,
		},
	})
}
