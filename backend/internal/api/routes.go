package api

import (
	"net/http"

	"github.com/dukerupert/coffee-commerce/internal/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, productHandler handler.ProductHandler, variantHandler handler.VariantHandler, priceHandler handler.PriceHandler, stripeWebhookHandler handler.StripeWebhookHandler, adminHandler handler.AdminHandler) error {

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	api := e.Group("/api")
	v1 := api.Group("/v1")
	v1.POST("/webhooks/stripe", stripeWebhookHandler.HandleWebhook)

	// Existing routes
	products := v1.Group("/products")
	products.GET("", productHandler.List)
	products.POST("", productHandler.Create)
	products.GET("/:id", productHandler.Get)
	products.DELETE("/:id", productHandler.Delete)
	products.GET("/:id/variants", variantHandler.ListByProduct)
	products.POST("/:id/archive", productHandler.Archive)

	// Add product prices route
	products.GET("/:id/prices", priceHandler.GetByProduct)

	// Add price routes
	prices := v1.Group("/prices")
	prices.POST("", priceHandler.Create)
	prices.GET("/:id", priceHandler.Get)
	prices.PUT("/:id", priceHandler.Update)
	prices.DELETE("/:id", priceHandler.Delete)
	prices.GET("/:id/variants", priceHandler.GetVariantsByPrice)

	// Add variant price assignment route
	variants := v1.Group("/variants")
	variants.POST("/:id/assign-price", priceHandler.AssignToVariant)

	// Admin routes (add these)
	admin := v1.Group("/admin")
	admin.GET("/health", adminHandler.HealthCheck)
	admin.POST("/sync-stripe-ids", adminHandler.SyncStripeProductIDs)

	return nil
}
