// internal/service/product_service.go
package service

import (
	"context"
	"fmt"

	"github.com/dukerupert/coffee-commerce/internal/domain/dto"
	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	events "github.com/dukerupert/coffee-commerce/internal/event"
	interfaces "github.com/dukerupert/coffee-commerce/internal/repository/interface"
	"github.com/rs/zerolog"
)

// ProductService defines the interface for product-related operations
type ProductService interface {
	// Create(ctx context.Context, productDTO *dto.ProductCreateDTO) (*models.Product, error)
	Create(ctx context.Context, product *dto.ProductCreateDTO) (*model.Product, error)
	// GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	List(ctx context.Context, offset, limit int, includeInactive bool) ([]*model.Product, int, error)
	// Update(ctx context.Context, id uuid.UUID, productDTO *dto.ProductUpdateDTO) (*models.Product, error)
	// Delete(ctx context.Context, id uuid.UUID) error
	// UpdateStockLevel(ctx context.Context, id uuid.UUID, quantity int) error
}

// productService implements ProductService
type productService struct {
	logger   zerolog.Logger
	eventBus events.EventBus
	repo     interfaces.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(logger *zerolog.Logger, eventBus events.EventBus, productRepo interfaces.ProductRepository) ProductService {
	subLogger := logger.With().Str("component", "product_service").Logger()
	return &productService{
		logger:   subLogger,
		eventBus: eventBus,
		repo:     productRepo,
	}
}

// Create initiates the product creation flow by saving to database and publishing an event
func (s *productService) Create(ctx context.Context, p *dto.ProductCreateDTO) (*model.Product, error) {
	s.logger.Info().
		Str("product_name", p.Name).
		Str("origin", p.Origin).
		Str("roast_level", p.RoastLevel).
		Msg("Creating product")

	// Convert dto to model
	product := p.ToModel()

	// Check if a product with the same name already exists
	existingProduct, err := s.repo.GetByName(ctx, product.Name)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error checking for existing product")
		return product, fmt.Errorf("error checking for existing product: %w", err)
	}

	if existingProduct != nil {
		s.logger.Warn().
			Str("product_name", product.Name).
			Str("existing_id", existingProduct.ID.String()).
			Msg("Product with this name already exists")
		return product, fmt.Errorf("a product with the name '%s' already exists", product.Name)
	}

	// Save product to database using repository
	err = s.repo.Create(ctx, product)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to save product to database")
		return product, fmt.Errorf("failed to create product: %w", err)
	}

	// Create event payload with important product details
	payload := events.ProductCreatedPayload{
		ProductID:         product.ID.String(),
		Name:              product.Name,
		Description:       product.Description,
		ImageURL:          product.ImageURL,
		StockLevel:        product.StockLevel,
		Weight:            product.Weight,
		Origin:            product.Origin,
		RoastLevel:        product.RoastLevel,
		FlavorNotes:       product.FlavorNotes,
		Options:           product.Options,
		AllowSubscription: product.AllowSubscription,
		Active:            product.Active,
		CreatedAt:         product.CreatedAt,
	}

	// Publish product created event with detailed payload
	err = s.eventBus.Publish(events.TopicProductCreated, payload)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to publish product created event")
		// Note: We don't return error here since the product was already saved to DB
		// Consider adding a mechanism to retry publishing the event
	}

	s.logger.Info().
		Str("topic", events.TopicProductCreated).
		Str("product_id", product.ID.String()).
		Msg("Published product created event")

	return product, nil
}

func (s *productService) List(ctx context.Context, offset, limit int, includeInactive bool) ([]*model.Product, int, error) {
	s.logger.Debug().
		Str("function", "productService.List").
		Int("offset", offset).
		Int("limit", limit).
		Bool("includeInactive", includeInactive).
		Msg("Starting product listing")

	products, total, err := s.repo.List(ctx, offset, limit, includeInactive)
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
