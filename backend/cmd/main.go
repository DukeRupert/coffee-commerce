package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/dukerupert/coffee-commerce/config"
	"github.com/dukerupert/coffee-commerce/internal/api/handler"
	events "github.com/dukerupert/coffee-commerce/internal/event"
	"github.com/dukerupert/coffee-commerce/internal/metrics"
	custommiddleware "github.com/dukerupert/coffee-commerce/internal/middleware"
	"github.com/dukerupert/coffee-commerce/internal/repository/postgres"
	"github.com/dukerupert/coffee-commerce/internal/service"
	"github.com/dukerupert/coffee-commerce/internal/stripe"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

func init() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	// Initialize logger
	logger  := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnix}).With().Timestamp().Logger()

	debug := flag.Bool("debug", false, "sets log level to debug")
	metricsAddr := flag.String("metrics-addr", ":9090", "The address the metrics server binds to")

	flag.Parse()
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database
	db, err := postgres.Connect(cfg.DB, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(cfg, &logger); err != nil {
		logger.Fatal().Err(err).Msg("Fatal migration error")
	}

	// Initialize metrics
	eventMetrics := metrics.NewEventMetrics()

	// Start metrics server
	go func() {
		logger.Info().Str("addr", *metricsAddr).Msg("Starting metrics server")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(*metricsAddr, nil); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start metrics server")
		}
	}()

	// Initialize event bus
	logger.Info().Msg("Initializing event bus")
	eventBus, err := events.NewNATSEventBus(
		cfg.MessageBus.URL,
		&logger,
		eventMetrics,
		"main-service",
	)
	defer eventBus.Close()

	// Initialize repositories
	productRepo := postgres.NewProductRepository(db, &logger)
	variantRepo := postgres.NewVariantRepository(db, &logger)
	priceRepo := postgres.NewPriceRepository(db, &logger)

	// Initialize services
	stripeService := stripe.NewStripeService(&logger, &cfg.Stripe)
	productService := service.NewProductService(&logger, eventBus, productRepo)
	_, err = service.NewVariantService(&logger, eventBus, variantRepo, productRepo, priceRepo, stripeService)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize variant service")
	}

	// Initialize handlers
	productHandler := handler.NewProductHandler(&logger, productService, variantRepo, priceRepo)
	variantHandler := handler.NewVariantHandler(&logger, variantRepo, productRepo)
	stripeWebhookHandler := handler.NewStripeWebhookHandler(&logger, &cfg.Stripe, eventBus)

	// Start echo server
	e := echo.New()

	// middleware
	e.Use(middleware.RequestID())
	e.Use(custommiddleware.RequestLogger(&logger))
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

	api := e.Group("/api")
	v1 := api.Group("/v1")
	v1.POST("/webhooks/stripe", stripeWebhookHandler.HandleWebhook)
	products := v1.Group("/products")
	products.GET("/", productHandler.List)
	products.POST("/", productHandler.Create)
	products.GET("/:id", productHandler.Get)
	products.DELETE("/:id", productHandler.Delete)
	products.GET("/:id/variants", variantHandler.ListByProduct)
	products.POST("/:id/archive", productHandler.Archive)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func runMigrations(cfg *config.Config, logger *zerolog.Logger) error {
	// Run migrations
	m, err := migrate.New(
		"file://migrations",
		cfg.DB.MigrateURL)
	if err != nil {
		logger.Error().Err(err).Str("migrateURL", cfg.DB.MigrateURL).Msg("Failed to create migration instance")
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Log migration source and database URL
	logger.Debug().
		Str("source", "file://migrations").
		Str("migrateURL", cfg.DB.MigrateURL).
		Msg("Migration configuration")

	// Get migration version before running
	version, dirty, vErr := m.Version()
	if vErr != nil && vErr != migrate.ErrNilVersion {
		logger.Warn().Err(vErr).Msg("Failed to get current migration version")
	} else if vErr == migrate.ErrNilVersion {
		logger.Info().Msg("No migrations have been applied yet")
	} else {
		logger.Info().Uint("version", version).Bool("dirty", dirty).Msg("Current migration version")
	}

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info().Msg("No migration changes detected")
		} else {
			logger.Error().Err(err).Msg("Migration failed")
			return fmt.Errorf("migration failed: %w", err)
		}
	} else {
		// Get the new version after successful migration
		newVersion, _, _ := m.Version()
		logger.Info().Uint("new_version", newVersion).Msg("Database migrations completed successfully")
	}

	// Close the migration
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		logger.Warn().Err(srcErr).Msg("Error closing migration source")
	}
	if dbErr != nil {
		logger.Warn().Err(dbErr).Msg("Error closing migration database connection")
	}

	// If both closing errors occurred, return a combined error
	if srcErr != nil && dbErr != nil {
		return fmt.Errorf("failed to close migration resources: %v, %v", srcErr, dbErr)
	} else if srcErr != nil {
		return fmt.Errorf("failed to close migration source: %w", srcErr)
	} else if dbErr != nil {
		return fmt.Errorf("failed to close migration database connection: %w", dbErr)
	}

	return nil
}
