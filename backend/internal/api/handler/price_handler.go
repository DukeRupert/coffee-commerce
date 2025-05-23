// internal/api/handler/price_handler.go
package handler

import (
	"errors"
	"net/http"

	"github.com/dukerupert/coffee-commerce/internal/api"
	"github.com/dukerupert/coffee-commerce/internal/domain/dto"
	"github.com/dukerupert/coffee-commerce/internal/interfaces"
	"github.com/dukerupert/coffee-commerce/internal/repository/postgres"
	"github.com/dukerupert/coffee-commerce/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type PriceHandler interface {
	Create(c echo.Context) error
	Get(c echo.Context) error
	GetByProduct(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	AssignToVariant(c echo.Context) error
	GetVariantsByPrice(c echo.Context) error
}

// priceHandler handles HTTP requests for prices
type priceHandler struct {
	logger       zerolog.Logger
	priceService interfaces.PriceService
	productRepo  interfaces.ProductRepository
	variantRepo  interfaces.VariantRepository
}

// NewPriceHandler creates a new price handler
func NewPriceHandler(
	logger *zerolog.Logger,
	priceService interfaces.PriceService,
	productRepo interfaces.ProductRepository,
	variantRepo interfaces.VariantRepository,
) *priceHandler {
	sublogger := logger.With().Str("component", "price_handler").Logger()
	return &priceHandler{
		logger:       sublogger,
		priceService: priceService,
		productRepo:  productRepo,
		variantRepo:  variantRepo,
	}
}

// Create handles POST /api/prices
func (h *priceHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	h.logger.Info().
		Str("handler", "PriceHandler.Create").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Msg("Handling price creation request")

	// Parse the request body
	var priceDTO dto.PriceCreateDTO
	if err := c.Bind(&priceDTO); err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Msg("Failed to parse request body")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
			Code:    "INVALID_FORMAT",
		})
	}

	// Validate the DTO
	validationErrors := priceDTO.Valid(ctx)
	if len(validationErrors) > 0 {
		h.logger.Warn().
			Interface("validation_errors", validationErrors).
			Str("request_id", requestID).
			Msg("Price validation failed")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:           http.StatusBadRequest,
			Message:          "Validation failed",
			ValidationErrors: validationErrors,
			Code:             "VALIDATION_ERROR",
		})
	}

	// Create the price
	price, err := h.priceService.Create(ctx, &priceDTO)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("product_id", priceDTO.ProductID.String()).
			Msg("Failed to create price")

		// Handle specific error types
		switch {
		case errors.Is(err, postgres.ErrResourceNotFound):
			return c.JSON(http.StatusNotFound, api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Product not found",
				Code:    "PRODUCT_NOT_FOUND",
			})

		case errors.Is(err, postgres.ErrDatabaseConnection):
			return c.JSON(http.StatusServiceUnavailable, api.ErrorResponse{
				Status:  http.StatusServiceUnavailable,
				Message: "Service temporarily unavailable, please try again later",
				Code:    "SERVICE_UNAVAILABLE",
			})

		case errors.Is(err, service.ErrInsufficientPermissions):
			return c.JSON(http.StatusForbidden, api.ErrorResponse{
				Status:  http.StatusForbidden,
				Message: "You don't have permission to create prices",
				Code:    "FORBIDDEN",
			})

		default:
			return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to create price",
				Code:    "INTERNAL_ERROR",
			})
		}
	}

	// Convert to response DTO
	priceResponse := dto.PriceResponseDTOFromModel(price)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Price created successfully",
		"price":   priceResponse,
	})
}

// Get handles GET /api/prices/:id
func (h *priceHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// Parse ID from URL
	idParam := c.Param("id")

	h.logger.Info().
		Str("handler", "PriceHandler.Get").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Str("id_param", idParam).
		Msg("Handling get price by ID request")

	// Convert string ID to UUID
	priceID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Str("id_param", idParam).
			Msg("Invalid price ID format")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid price ID format",
			Code:    "INVALID_ID_FORMAT",
		})
	}

	// Get price from service
	price, err := h.priceService.GetByID(ctx, priceID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("price_id", priceID.String()).
			Msg("Failed to retrieve price")

		if errors.Is(err, postgres.ErrResourceNotFound) {
			return c.JSON(http.StatusNotFound, api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Price not found",
				Code:    "PRICE_NOT_FOUND",
			})
		}

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve price",
			Code:    "INTERNAL_ERROR",
		})
	}

	// Convert to response DTO
	priceResponse := dto.PriceResponseDTOFromModel(price)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"price": priceResponse,
	})
}

// GetByProduct handles GET /api/products/:id/prices
func (h *priceHandler) GetByProduct(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// Parse product ID from URL
	productIDParam := c.Param("id")

	h.logger.Info().
		Str("handler", "PriceHandler.GetByProduct").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Str("product_id", productIDParam).
		Msg("Handling get prices by product ID request")

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

	// Get prices for the product
	prices, err := h.priceService.GetByProductID(ctx, productID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("product_id", productID.String()).
			Msg("Failed to retrieve prices for product")

		if errors.Is(err, postgres.ErrResourceNotFound) {
			return c.JSON(http.StatusNotFound, api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Product not found",
				Code:    "PRODUCT_NOT_FOUND",
			})
		}

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve prices",
			Code:    "INTERNAL_ERROR",
		})
	}

	// Convert to response DTOs
	priceResponses := make([]*dto.PriceResponseDTO, len(prices))
	for i, price := range prices {
		priceResponse := dto.PriceResponseDTOFromModel(price)
		priceResponses[i] = &priceResponse
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"prices": priceResponses,
		"count":  len(priceResponses),
	})
}

// Update handles PUT /api/prices/:id
func (h *priceHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// Parse ID from URL
	idParam := c.Param("id")

	h.logger.Info().
		Str("handler", "PriceHandler.Update").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Str("id_param", idParam).
		Msg("Handling price update request")

	// Convert string ID to UUID
	priceID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Str("id_param", idParam).
			Msg("Invalid price ID format")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid price ID format",
			Code:    "INVALID_ID_FORMAT",
		})
	}

	// Parse the request body
	var updateDTO dto.PriceUpdateDTO
	if err := c.Bind(&updateDTO); err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Msg("Failed to parse request body")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
			Code:    "INVALID_FORMAT",
		})
	}

	// Validate the DTO
	validationErrors := updateDTO.Valid(ctx)
	if len(validationErrors) > 0 {
		h.logger.Warn().
			Interface("validation_errors", validationErrors).
			Str("request_id", requestID).
			Msg("Price update validation failed")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:           http.StatusBadRequest,
			Message:          "Validation failed",
			ValidationErrors: validationErrors,
			Code:             "VALIDATION_ERROR",
		})
	}

	// Update the price
	price, err := h.priceService.Update(ctx, priceID, &updateDTO)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("price_id", priceID.String()).
			Msg("Failed to update price")

		if errors.Is(err, postgres.ErrResourceNotFound) {
			return c.JSON(http.StatusNotFound, api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Price not found",
				Code:    "PRICE_NOT_FOUND",
			})
		}

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update price",
			Code:    "INTERNAL_ERROR",
		})
	}

	// Convert to response DTO
	priceResponse := dto.PriceResponseDTOFromModel(price)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Price updated successfully",
		"price":   priceResponse,
	})
}

// Delete handles DELETE /api/prices/:id
func (h *priceHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// Parse ID from URL
	idParam := c.Param("id")

	h.logger.Info().
		Str("handler", "PriceHandler.Delete").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Str("id_param", idParam).
		Msg("Handling price deletion request")

	// Convert string ID to UUID
	priceID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Str("id_param", idParam).
			Msg("Invalid price ID format")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid price ID format",
			Code:    "INVALID_ID_FORMAT",
		})
	}

	// Delete the price
	err = h.priceService.Delete(ctx, priceID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("price_id", priceID.String()).
			Msg("Failed to delete price")

		if errors.Is(err, postgres.ErrResourceNotFound) {
			return c.JSON(http.StatusNotFound, api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Price not found",
				Code:    "PRICE_NOT_FOUND",
			})
		}

		// Check for foreign key constraint errors (variants using this price)
		if err.Error() == "cannot delete price: 1 variants are using this price" ||
			err.Error() == "cannot delete price: 2 variants are using this price" {
			return c.JSON(http.StatusConflict, api.ErrorResponse{
				Status:  http.StatusConflict,
				Message: "Cannot delete price because it is being used by variants",
				Code:    "PRICE_IN_USE",
			})
		}

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete price",
			Code:    "INTERNAL_ERROR",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Price deleted successfully",
	})
}

// AssignToVariant handles POST /api/variants/:id/assign-price
func (h *priceHandler) AssignToVariant(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// Parse variant ID from URL
	variantIDParam := c.Param("id")

	h.logger.Info().
		Str("handler", "PriceHandler.AssignToVariant").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Str("variant_id", variantIDParam).
		Msg("Handling assign price to variant request")

	// Convert string ID to UUID
	variantID, err := uuid.Parse(variantIDParam)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Str("variant_id", variantIDParam).
			Msg("Invalid variant ID format")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid variant ID format",
			Code:    "INVALID_ID_FORMAT",
		})
	}

	// Parse the request body to get the price ID
	var requestBody struct {
		PriceID uuid.UUID `json:"price_id"`
	}
	if err := c.Bind(&requestBody); err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Msg("Failed to parse request body")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
			Code:    "INVALID_FORMAT",
		})
	}

	// Create assignment DTO
	assignmentDTO := &dto.VariantPriceAssignmentDTO{
		VariantID: variantID,
		PriceID:   requestBody.PriceID,
	}

	// Validate the DTO
	validationErrors := assignmentDTO.Valid(ctx)
	if len(validationErrors) > 0 {
		h.logger.Warn().
			Interface("validation_errors", validationErrors).
			Str("request_id", requestID).
			Msg("Price assignment validation failed")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:           http.StatusBadRequest,
			Message:          "Validation failed",
			ValidationErrors: validationErrors,
			Code:             "VALIDATION_ERROR",
		})
	}

	// Assign the price to the variant
	err = h.priceService.AssignToVariant(ctx, assignmentDTO)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("variant_id", variantID.String()).
			Str("price_id", requestBody.PriceID.String()).
			Msg("Failed to assign price to variant")

		if errors.Is(err, postgres.ErrResourceNotFound) {
			return c.JSON(http.StatusNotFound, api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Variant or price not found",
				Code:    "RESOURCE_NOT_FOUND",
			})
		}

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to assign price to variant",
			Code:    "INTERNAL_ERROR",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Price assigned to variant successfully",
	})
}

// GetVariantsByPrice handles GET /api/prices/:id/variants
func (h *priceHandler) GetVariantsByPrice(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// Parse price ID from URL
	priceIDParam := c.Param("id")

	h.logger.Info().
		Str("handler", "PriceHandler.GetVariantsByPrice").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Str("price_id", priceIDParam).
		Msg("Handling get variants by price ID request")

	// Convert string ID to UUID
	priceID, err := uuid.Parse(priceIDParam)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Str("price_id", priceIDParam).
			Msg("Invalid price ID format")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid price ID format",
			Code:    "INVALID_ID_FORMAT",
		})
	}

	// Get variants using this price
	variants, err := h.priceService.GetVariantsByPrice(ctx, priceID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("price_id", priceID.String()).
			Msg("Failed to retrieve variants for price")

		if errors.Is(err, postgres.ErrResourceNotFound) {
			return c.JSON(http.StatusNotFound, api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Price not found",
				Code:    "PRICE_NOT_FOUND",
			})
		}

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve variants",
			Code:    "INTERNAL_ERROR",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"variants": variants,
		"count":    len(variants),
	})
}