package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dukerupert/coffee-commerce/internal/api"
	"github.com/dukerupert/coffee-commerce/internal/domain/dto"
	"github.com/dukerupert/coffee-commerce/internal/interfaces"
	"github.com/dukerupert/coffee-commerce/internal/repository/postgres"
	"github.com/dukerupert/coffee-commerce/internal/service"
	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ProductHandler interface {
	Create(c echo.Context) error
	Get(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
	Archive(c echo.Context) error
	Delete(c echo.Context) error
	UpdateStockLevel(c echo.Context) error
}

// ProductHandler handles HTTP requests for products
type productHandler struct {
	logger         zerolog.Logger
	productService interfaces.ProductService
	variantRepo    interfaces.VariantRepository
	priceRepo      interfaces.PriceRepository
}

// NewProductHandler creates a new product handler
func NewProductHandler(logger *zerolog.Logger, productService interfaces.ProductService, variantRepo interfaces.VariantRepository, priceRepo interfaces.PriceRepository) *productHandler {
	sublogger := logger.With().Str("component", "product_handler").Logger()
	return &productHandler{
		logger:         sublogger,
		productService: productService,
		variantRepo:    variantRepo,
		priceRepo:      priceRepo,
	}
}

// Create handles POST /api/products
func (h *productHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	h.logger.Info().
		Str("handler", "ProductHandler.Create").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Msg("Handling product creation request")

	// Parse the request body
	var productDTO dto.ProductCreateDTO
	if err := c.Bind(&productDTO); err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Msg("Failed to parse request body")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "Invalid request format",
			Code:    "INVALID_FORMAT",
		})
	}

	// Set default weight
	if productDTO.Weight < 1 {
		productDTO.Weight = 1
	}

	// Validate using the existing Valid method
	validationErrors := productDTO.Valid(c.Request().Context())
	if len(validationErrors) > 0 {
		h.logger.Warn().
			Interface("validation_errors", validationErrors).
			Str("request_id", requestID).
			Msg("Product validation failed")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message:          "Validation failed",
			ValidationErrors: validationErrors,
			Code:             "VALIDATION_ERROR",
		})
	}

	// Call product service to create the product
	product, err := h.productService.Create(ctx, &productDTO)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Msg("Failed to create product")

		// Handle specific error types
		switch {
		case errors.Is(err, postgres.ErrDuplicateName):
			return c.JSON(http.StatusConflict, api.ErrorResponse{
				Message: "A product with this name already exists",
				Code:    "DUPLICATE_PRODUCT",
				ValidationErrors: map[string]string{
					"name": "This product name is already in use",
				},
			})

		case errors.Is(err, postgres.ErrDatabaseConnection):
			return c.JSON(http.StatusServiceUnavailable, api.ErrorResponse{
				Message: "Service temporarily unavailable, please try again later",
				Code:    "SERVICE_UNAVAILABLE",
			})

		case errors.Is(err, service.ErrInsufficientPermissions):
			return c.JSON(http.StatusForbidden, api.ErrorResponse{
				Message: "You don't have permission to create products",
				Code:    "FORBIDDEN",
			})

		default:
			// Generic server error
			return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Message: "Failed to create product",
				Code:    "INTERNAL_ERROR",
			})
		}
	}

	// Return success response with the created product ID
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Product created successfully",
		"product": product,
	})
}

// Get handles GET /api/products/:id
func (h *productHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// 1. Parse ID from URL
	idParam := c.Param("id")

	h.logger.Info().
		Str("handler", "ProductHandler.Get").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Str("id_param", idParam).
		Msg("Handling get product by ID request")

		// Convert string ID to UUID
	productID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Str("id_param", idParam).
			Msg("Invalid product ID format")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid product ID format",
			Code:    "INVALID_ID_FORMAT",
		})
	}

	// 2. Get product from database
	product, err := h.productService.GetByID(ctx, productID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("product_id", productID.String()).
			Msg("Failed to retrieve product")

		if errors.Is(err, postgres.ErrResourceNotFound) {
			return c.JSON(http.StatusNotFound, api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Product not found",
				Code:    "PRODUCT_NOT_FOUND",
			})
		}

		return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve product",
			Code:    "INTERNAL_ERROR",
		})
	}

	// 3. Convert product to DTO
	productResponse := dto.ProductResponseDTOFromModel(product)

	// 4. Get variants for this product
	variants, err := h.variantRepo.GetByProductID(ctx, productID)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("product_id", productID.String()).
			Msg("Failed to retrieve variants, continuing with product only")
		// Continue with just the product info if variant retrieval fails
	}

	// 5. Build variant responses with prices
	variantResponses := []map[string]interface{}{}

	for _, variant := range variants {
		// Get price information for this variant
		price, err := h.priceRepo.GetByID(ctx, variant.PriceID)
		if err != nil {
			h.logger.Warn().
				Err(err).
				Str("variant_id", variant.ID.String()).
				Str("price_id", variant.PriceID.String()).
				Msg("Failed to retrieve price for variant")
			// Continue with basic variant info if price retrieval fails
		}

		variantResponse := map[string]interface{}{
			"id":              variant.ID.String(),
			"product_id":      variant.ProductID.String(),
			"price_id":        variant.PriceID.String(),
			"stripe_price_id": variant.StripePriceID,
			"options":         variant.Options,
			"active":          variant.Active,
			"stock_level":     variant.StockLevel,
			"created_at":      variant.CreatedAt.Format(time.RFC3339),
			"updated_at":      variant.UpdatedAt.Format(time.RFC3339),
		}

		// Add price information if available
		if price != nil {
			variantResponse["price"] = map[string]interface{}{
				"amount":   price.Amount,
				"currency": price.Currency,
				"type":     price.Type,
			}
		}

		variantResponses = append(variantResponses, variantResponse)
	}

	// 6. Include variants in the response
	response := map[string]interface{}{
		"product":  productResponse,
		"variants": variantResponses,
	}

	return c.JSON(http.StatusOK, response)
}

// List handles GET /api/products
func (h *productHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	h.logger.Debug().
		Str("handler", "ProductHandler.List").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Msg("Handling product listing request")

	// 1. Parse pagination parameters
	params := api.NewParams(c)

	// 2. Parse additional filtering parameters
	includeInactive := false
	if c.QueryParam("include_inactive") == "true" {
		// todo: Only admins to see inactive products
		includeInactive = true
		h.logger.Debug().
			Str("handler", "ProductHandler.List").
			Str("request_id", requestID).
			Bool("include_inactive", includeInactive).
			Msg("Including inactive products in results")
	}
	
	includeArchived := false
	if c.QueryParam("include_archived") == "true" {
		// todo: Only admins to see archived products
		includeArchived = true
		h.logger.Debug().
			Str("handler", "ProductHandler.List").
			Str("request_id", requestID).
			Bool("include_archived", includeArchived).
			Msg("Including archived products in results")
	}

	// 2. Call productService.List
	products, total, err := h.productService.List(ctx, params.Offset, params.PerPage, includeInactive, includeArchived)
	if err != nil {
		h.logger.Error().
			Str("handler", "ProductHandler.List").
			Str("request_id", requestID).
			Err(err).
			Int("offset", params.Offset).
			Int("per_page", params.PerPage).
			Bool("include_inactive", includeInactive).
			Bool("include_archived", includeArchived).
			Msg("Failed to retrieve products from service")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve products")
	}

	// 3. Create paginated response
	meta := api.NewMeta(params, total)
	response := api.Response(products, meta)

	h.logger.Info().
		Str("handler", "ProductHandler.List").
		Str("request_id", requestID).
		Int("products_count", len(products)).
		Int("total_count", total).
		Int("page", params.Page).
		Int("per_page", params.PerPage).
		Int("status_code", http.StatusOK).
		Msg("Product listing successfully returned")

	return c.JSON(http.StatusOK, response)
}

// Update handles PUT /api/products/:id
func (h *productHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// 1. Parse ID from URL
	idParam := c.Param("id")

	h.logger.Info().
		Str("handler", "ProductHandler.Update").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Str("id_param", idParam).
		Msg("Handling product update by ID request")

	// Convert string ID to UUID
	productID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Str("id_param", idParam).
			Msg("Invalid product ID format")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid product ID format",
			Code:    "INVALID_ID_FORMAT",
		})
	}

	// 2. Parse the request body
	var productUpdateDTO dto.ProductUpdateDTO
	if err := c.Bind(&productUpdateDTO); err != nil {
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

	// 3. Validate the DTO
	validationErrors := productUpdateDTO.Valid(ctx)
	if len(validationErrors) > 0 {
		h.logger.Warn().
			Interface("validation_errors", validationErrors).
			Str("request_id", requestID).
			Msg("Product update validation failed")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:           http.StatusBadRequest,
			Message:          "Validation failed",
			ValidationErrors: validationErrors,
			Code:             "VALIDATION_ERROR",
		})
	}

	// 4. Call product service to update the product
	updatedProduct, err := h.productService.Update(ctx, productID, &productUpdateDTO)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("product_id", productID.String()).
			Msg("Failed to update product")

		// Handle specific error types
		switch {
		case errors.Is(err, postgres.ErrResourceNotFound):
			return c.JSON(http.StatusNotFound, api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Product not found",
				Code:    "PRODUCT_NOT_FOUND",
			})

		case errors.Is(err, postgres.ErrDuplicateName):
			return c.JSON(http.StatusConflict, api.ErrorResponse{
				Status:  http.StatusConflict,
				Message: "A product with this name already exists",
				Code:    "DUPLICATE_PRODUCT",
				ValidationErrors: map[string]string{
					"name": "This product name is already in use",
				},
			})

		case errors.Is(err, service.ErrInsufficientPermissions):
			return c.JSON(http.StatusForbidden, api.ErrorResponse{
				Status:  http.StatusForbidden,
				Message: "You don't have permission to update this product",
				Code:    "FORBIDDEN",
			})

		case errors.Is(err, postgres.ErrDatabaseConnection):
			return c.JSON(http.StatusServiceUnavailable, api.ErrorResponse{
				Status:  http.StatusServiceUnavailable,
				Message: "Service temporarily unavailable, please try again later",
				Code:    "SERVICE_UNAVAILABLE",
			})

		default:
			// Generic server error
			return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to update product",
				Code:    "INTERNAL_ERROR",
			})
		}
	}

	// 5. Convert to response DTO
	productResponse := dto.ProductResponseDTOFromModel(updatedProduct)

	// 6. Return success response
	h.logger.Info().
		Str("request_id", requestID).
		Str("product_id", productID.String()).
		Str("product_name", updatedProduct.Name).
		Msg("Product updated successfully")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product updated successfully",
		"product": productResponse,
	})
}

// Archive handles POST /api/products/:id/archive
func (h *productHandler) Archive(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	// 1. Parse ID from URL
	idParam := c.Param("id")

	h.logger.Info().
		Str("handler", "ProductHandler.Archive").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Str("id_param", idParam).
		Msg("Handling product archive by ID request")

	// Convert string ID to UUID
	productID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Str("id_param", idParam).
			Msg("Invalid product ID format")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid product ID format",
			Code:    "INVALID_ID_FORMAT",
		})
	}

	// 2. Call service to archive the product
	err = h.productService.Archive(ctx, productID)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("product_id", productID.String()).
			Msg("Failed to archive product")

		// Handle specific error types
		switch {
		case errors.Is(err, postgres.ErrResourceNotFound):
			return c.JSON(http.StatusNotFound, api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Product not found",
				Code:    "PRODUCT_NOT_FOUND",
			})

		case errors.Is(err, service.ErrInsufficientPermissions):
			return c.JSON(http.StatusForbidden, api.ErrorResponse{
				Status:  http.StatusForbidden,
				Message: "You don't have permission to archive this product",
				Code:    "FORBIDDEN",
			})

		default:
			// Generic server error
			return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to archive product",
				Code:    "INTERNAL_ERROR",
			})
		}
	}

	// 3. Return success response
	h.logger.Info().
		Str("request_id", requestID).
		Str("product_id", productID.String()).
		Msg("Product archived successfully")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Product archived successfully",
	})
}

// Delete handles DELETE /api/products/:id
// Delete now redirects to Archive as the preferred method for removing products
func (h *productHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	
	// Get the "hard_delete" query parameter
	hardDelete := c.QueryParam("hard_delete") == "true"
	
	// 1. Parse ID from URL
	idParam := c.Param("id")

	h.logger.Info().
		Str("handler", "ProductHandler.Delete").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Str("id_param", idParam).
		Bool("hard_delete", hardDelete).
		Msg("Handling product delete by ID request")

	// Convert string ID to UUID
	productID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("request_id", requestID).
			Str("id_param", idParam).
			Msg("Invalid product ID format")

		return c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid product ID format",
			Code:    "INVALID_ID_FORMAT",
		})
	}

	// If not a hard delete request, redirect to Archive method
	if !hardDelete {
		h.logger.Info().
			Str("request_id", requestID).
			Str("product_id", productID.String()).
			Msg("Redirecting delete request to archive operation")
			
		err = h.productService.Archive(ctx, productID)
	} else {
		// This is a hard delete request - only allow for admin users
		// TODO: Add proper permission checking here
		h.logger.Warn().
			Str("request_id", requestID).
			Str("product_id", productID.String()).
			Msg("Attempting hard delete operation - this should only be used for testing")
			
		// Attempt the hard delete
		err = h.productService.Delete(ctx, productID)
	}
	
	// Error handling - same for both archive and delete
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("request_id", requestID).
			Str("product_id", productID.String()).
			Bool("hard_delete", hardDelete).
			Msg("Failed to remove product")

		// Handle specific error types
		switch {
		case errors.Is(err, postgres.ErrResourceNotFound):
			return c.JSON(http.StatusNotFound, api.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Product not found",
				Code:    "PRODUCT_NOT_FOUND",
			})

		case errors.Is(err, service.ErrInsufficientPermissions):
			return c.JSON(http.StatusForbidden, api.ErrorResponse{
				Status:  http.StatusForbidden,
				Message: "You don't have permission to delete this product",
				Code:    "FORBIDDEN",
			})
			
		case errors.Is(err, postgres.ErrDatabaseConnection):
			return c.JSON(http.StatusServiceUnavailable, api.ErrorResponse{
				Status:  http.StatusServiceUnavailable,
				Message: "Service temporarily unavailable, please try again later",
				Code:    "SERVICE_UNAVAILABLE",
			})

		default:
			// If it's a foreign key constraint error (common with hard delete), 
			// suggest using archive instead
			if hardDelete {
				return c.JSON(http.StatusConflict, api.ErrorResponse{
					Status:  http.StatusConflict,
					Message: "This product cannot be hard deleted because it has associated records. Use archive instead.",
					Code:    "FOREIGN_KEY_CONSTRAINT",
				})
			}
			
			// Generic server error
			return c.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to remove product",
				Code:    "INTERNAL_ERROR",
			})
		}
	}

	// 3. Return success response
	operation := "archived"
	if hardDelete {
		operation = "deleted"
	}
	
	h.logger.Info().
		Str("request_id", requestID).
		Str("product_id", productID.String()).
		Bool("hard_delete", hardDelete).
		Str("operation", operation).
		Msg("Product successfully removed")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Product successfully %s", operation),
	})
}

// UpdateStockLevel handles PATCH /api/products/:id/stock
func (h *productHandler) UpdateStockLevel(c echo.Context) error {
	h.logger.Info().Str("handler", "ProductHandler.UpdateStockLevel").Msg("Handling product update stock by ID request")
	return c.String(http.StatusOK, "Hello, World!")
}
