// internal/service/product_service.go
package service

import (
	"github.com/dukerupert/coffee-commerce/internal/events"
	"github.com/rs/zerolog"
)

// ProductService defines the interface for product-related operations
type ProductService interface {
	Create() error
	// Other methods will be added later
}

// productService implements ProductService
type productService struct {
	logger   zerolog.Logger
	eventBus events.EventBus
}

// NewProductService creates a new product service
func NewProductService(logger *zerolog.Logger, eventBus events.EventBus) ProductService {
	subLogger := logger.With().Str("component", "product_service").Logger()
	return &productService{
		logger:   subLogger,
		eventBus: eventBus,
	}
}

// Create initiates the product creation flow by publishing an event
func (s *productService) Create() error {
	s.logger.Info().Msg("Creating product")
	
	// Publish product created event
	// Note: We're publishing an empty payload for now
	err := s.eventBus.Publish(events.TopicProductCreated, struct{}{})
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to publish product created event")
		return err
	}
	
	s.logger.Info().Str("topic", events.TopicProductCreated).Msg("Published product created event")
	return nil
}