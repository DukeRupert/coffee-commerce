// internal/service/product_service.go
package service

import (
	"context"
	"fmt"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/dukerupert/coffee-commerce/internal/event"
	"github.com/dukerupert/coffee-commerce/internal/interfaces"
	"github.com/rs/zerolog"
)

// ProductService defines the interface for product-related operations
type ProductService interface {
	// Create(ctx context.Context, productDTO *dto.ProductCreateDTO) (*models.Product, error)
	Create() error
	// GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	List(ctx context.Context, offset, limit int, includeInactive bool) ([]*model.Product, int, error)
	// Update(ctx context.Context, id uuid.UUID, productDTO *dto.ProductUpdateDTO) (*models.Product, error)
	// Delete(ctx context.Context, id uuid.UUID) error
	// UpdateStockLevel(ctx context.Context, id uuid.UUID, quantity int) error
}

// productService implements ProductService
type productService struct {
	logger      zerolog.Logger
	eventBus    events.EventBus
	productRepo interfaces.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(logger *zerolog.Logger, eventBus events.EventBus, productRepo interfaces.ProductRepository) ProductService {
	subLogger := logger.With().Str("component", "product_service").Logger()
	return &productService{
		logger:      subLogger,
		eventBus:    eventBus,
		productRepo: productRepo,
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

func (s *productService) List(ctx context.Context, offset, limit int, includeInactive bool) ([]*model.Product, int, error) {
	s.logger.Debug().
		Str("function", "productService.List").
		Int("offset", offset).
		Int("limit", limit).
		Bool("includeInactive", includeInactive).
		Msg("Starting product listing")

	products, total, err := s.productRepo.List(ctx, offset, limit, includeInactive)
	if err != nil {
		s.logger.Error().
			Str("function", "productService.List").
			Err(err).
			Int("offset", offset).
			Int("limit", limit).
			Bool("includeInactive", includeInactive).
			Msg("Failed to retrieve products from repository")
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	s.logger.Info().
		Str("function", "productService.List").
		Int("total_products", total).
		Int("returned_products", len(products)).
		Int("offset", offset).
		Int("limit", limit).
		Bool("includeInactive", includeInactive).
		Msg("Product listing completed successfully")

	return products, total, nil
}
