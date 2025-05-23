package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/dukerupert/coffee-commerce/config"
	"github.com/dukerupert/coffee-commerce/internal/api"
	"github.com/dukerupert/coffee-commerce/internal/events"
	"github.com/dukerupert/coffee-commerce/internal/metrics"
	"github.com/dukerupert/coffee-commerce/internal/repository/postgres"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnix}).With().Timestamp().Logger()

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
	if err != nil {
		logger.Fatal().Err(err).Msg("Event bus failed to start")
	}
	defer eventBus.Close()

	// Debug
	logger.Debug().
		Str("stripe_key_status", func() string {
			if cfg.Stripe.SecretKey == "" {
				return "empty"
			}
			return "set"
		}()).
		Msg("Stripe configuration loaded")

	// Start echo server
	s := api.NewServer(cfg, db, eventBus, &logger)
	s.Start(cfg.App.Port)
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
