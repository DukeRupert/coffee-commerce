// internal/service/variant_service.go
package service

import (
	"encoding/json"
	
	"github.com/dukerupert/coffee-commerce/internal/events"
	"github.com/rs/zerolog"
)

// VariantService defines the interface for variant-related operations
type VariantService interface {
	// Methods will be added later
}

// variantService implements VariantService
type variantService struct {
	logger   zerolog.Logger
	eventBus events.EventBus
}

// NewVariantService creates a new variant service and subscribes to relevant events
func NewVariantService(logger *zerolog.Logger, eventBus events.EventBus) (VariantService, error) {
	subLogger := logger.With().Str("component", "variant_service").Logger()
	
	service := &variantService{
		logger:   subLogger,
		eventBus: eventBus,
	}
	
	// Subscribe to product created events
	_, err := eventBus.Subscribe(events.TopicProductCreated, service.handleProductCreated)
	if err != nil {
		subLogger.Error().Err(err).Msg("Failed to subscribe to product created events")
		return nil, err
	}
	
	subLogger.Info().Str("topic", events.TopicProductCreated).Msg("Subscribed to product created events")
	return service, nil
}

// handleProductCreated is called when a product created event is received
func (s *variantService) handleProductCreated(data []byte) {
	s.logger.Info().Str("topic", events.TopicProductCreated).Msg("Received product created event")
	
	// Parse the event
	var event events.Event
	if err := json.Unmarshal(data, &event); err != nil {
		s.logger.Error().Err(err).Msg("Failed to unmarshal product created event")
		return
	}
	
	// Process the event
	s.logger.Info().
		Str("event_id", event.ID).
		Time("timestamp", event.Timestamp).
		Msg("Processing product created event")
	
	// Publish variant created event
	err := s.eventBus.Publish(events.TopicVariantCreated, struct{}{})
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to publish variant created event")
		return
	}
	
	s.logger.Info().Str("topic", events.TopicVariantCreated).Msg("Published variant created event")
}