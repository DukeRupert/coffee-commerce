package api

import (
	"fmt"

	"github.com/dukerupert/coffee-commerce/config"
	"github.com/dukerupert/coffee-commerce/internal/events"
	"github.com/dukerupert/coffee-commerce/internal/handler"
	custommiddleware "github.com/dukerupert/coffee-commerce/internal/middleware"
	"github.com/dukerupert/coffee-commerce/internal/repository/postgres"
	"github.com/dukerupert/coffee-commerce/internal/service"
	"github.com/dukerupert/coffee-commerce/internal/stripe"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type server struct {
	e *echo.Echo
}

func NewServer(cfg *config.Config, db *postgres.DB, eventBus *events.NATSEventBus, logger *zerolog.Logger) *server {

	// Initialize repositories
	productRepo := postgres.NewProductRepository(db, logger)
	variantRepo := postgres.NewVariantRepository(db, logger)
	priceRepo := postgres.NewPriceRepository(db, logger)
	syncRepo := postgres.NewSyncHashRepository(db, logger)

	// Initialize services
	stripeService := stripe.NewStripeService(logger, &cfg.Stripe)
	productService := service.NewProductService(logger, eventBus, productRepo)
	priceService := service.NewPriceService(logger, eventBus, priceRepo, productRepo, variantRepo, stripeService)
	_, err := service.NewVariantService(logger, eventBus, variantRepo, productRepo, priceRepo, stripeService)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize variant service")
	}

	// Initialize handlers
	productHandler := handler.NewProductHandler(logger, productService, variantRepo, priceRepo)
	variantHandler := handler.NewVariantHandler(logger, variantRepo, productRepo)
	priceHandler := handler.NewPriceHandler(logger, priceService, productRepo, variantRepo)
	stripeWebhookHandler := handler.NewStripeWebhookHandler(logger, &cfg.Stripe, eventBus, productRepo, priceRepo, variantRepo, syncRepo)
	adminHandler := handler.NewAdminHandler(logger, priceService, productRepo)

	// Start echo server
	e := echo.New()

	// middleware
	e.Use(middleware.RequestID())
	e.Use(custommiddleware.RequestLogger(logger))
	corsConfig := custommiddleware.CORSConfig{
		AllowOrigins: []string{
			"https://orange-goldfish-wg644q6vqxv295gg-5173.app.github.dev", // Your frontend origin
			"http://localhost:5173", // Local development
			"*",                     // Allow all origins (for development)
		},
		AllowMethods: []string{
			echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS,
		},
	}
	e.Use(custommiddleware.SetupCORS(corsConfig))

	RegisterRoutes(e, productHandler, variantHandler, priceHandler, *stripeWebhookHandler, adminHandler)

	return &server{
		e: e,
	}
}

func (s *server) Start(port int) {
	s.e.Logger.Fatal(s.e.Start(fmt.Sprintf(":%d", port)))
}
