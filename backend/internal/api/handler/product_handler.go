package handler

import (
	"net/http"

	"github.com/dukerupert/coffee-commerce/internal/api"
	"github.com/dukerupert/coffee-commerce/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ProductHandler interface {
	Create(c echo.Context) error
	Get(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	UpdateStockLevel(c echo.Context) error
}

// ProductHandler handles HTTP requests for products
type productHandler struct {
	logger         zerolog.Logger
	productService service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(logger *zerolog.Logger, productService service.ProductService) *productHandler {
	sublogger := logger.With().Str("component", "product_handler").Logger()
	return &productHandler{
		logger:         sublogger,
		productService: productService,
	}
}

// Create handles POST /api/products
func (h *productHandler) Create(c echo.Context) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)

	h.logger.Info().
		Str("handler", "ProductHandler.Create").
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Str("remote_addr", c.Request().RemoteAddr).
		Msg("Handling product creation request")

	// Call product service to create the product
	err := h.productService.Create()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create product")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create product",
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Product creation initiated",
	})
}

// Get handles GET /api/products/:id
func (h *productHandler) Get(c echo.Context) error {
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
	return c.String(http.StatusOK, "Hello, World!")
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

	// 2. Call productService.List
	products, total, err := h.productService.List(ctx, params.Offset, params.PerPage, includeInactive)
	if err != nil {
		h.logger.Error().
			Str("handler", "ProductHandler.List").
			Str("request_id", requestID).
			Err(err).
			Int("offset", params.Offset).
			Int("per_page", params.PerPage).
			Bool("include_inactive", includeInactive).
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
	h.logger.Info().Str("handler", "ProductHandler.Update").Msg("Handling product update by ID request")
	return c.String(http.StatusOK, "Hello, World!")
}

// Delete handles DELETE /api/products/:id
func (h *productHandler) Delete(c echo.Context) error {
	h.logger.Info().Str("handler", "ProductHandler.Delete").Msg("Handling product delete by ID request")
	return c.String(http.StatusOK, "Hello, World!")
}

// UpdateStockLevel handles PATCH /api/products/:id/stock
func (h *productHandler) UpdateStockLevel(c echo.Context) error {
	h.logger.Info().Str("handler", "ProductHandler.UpdateStockLevel").Msg("Handling product update stock by ID request")
	return c.String(http.StatusOK, "Hello, World!")
}
