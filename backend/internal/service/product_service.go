// internal/service/product_service.go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dukerupert/coffee-commerce/internal/domain/dto"
	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/dukerupert/coffee-commerce/internal/interfaces"
	events "github.com/dukerupert/coffee-commerce/internal/event"
	"github.com/dukerupert/coffee-commerce/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// productService implements ProductService
type productService struct {
	logger   zerolog.Logger
	eventBus events.EventBus
	repo     interfaces.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(logger *zerolog.Logger, eventBus events.EventBus, productRepo interfaces.ProductRepository) interfaces.ProductService {
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

func (s *productService) List(ctx context.Context, offset, limit int, includeInactive, includeArchived bool) ([]*model.Product, int, error) {
	s.logger.Debug().
		Str("function", "productService.List").
		Int("offset", offset).
		Int("limit", limit).
		Bool("includeInactive", includeInactive).
		Bool("includeArchived", includeArchived).
		Msg("Starting product listing")

	products, total, err := s.repo.List(ctx, offset, limit, includeInactive, includeArchived)
	if err != nil {
		s.logger.Error().
			Str("function", "productService.List").
			Err(err).
			Int("offset", offset).
			Int("limit", limit).
			Bool("includeInactive", includeInactive).
			Bool("includeArchived", includeArchived).
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
		Bool("includeArchived", includeArchived).
		Msg("Product listing completed successfully")

	return products, total, nil
}

func (s *productService) GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	s.logger.Info().
		Str("product_id", id.String()).
		Msg("Retrieving product by ID")

		// Get product from repository
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("product_id", id.String()).
			Msg("Failed to retrieve product")
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}

	// Check if product exists
	if product == nil {
		s.logger.Warn().
			Str("product_id", id.String()).
			Msg("Product not found")
		return nil, postgres.ErrResourceNotFound
	}

	return product, nil
}

// Update updates an existing product
func (s *productService) Update(ctx context.Context, id uuid.UUID, dto *dto.ProductUpdateDTO) (*model.Product, error) {
	s.logger.Info().
		Str("product_id", id.String()).
		Msg("Updating product")

	// First, get the existing product
	existingProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("product_id", id.String()).
			Msg("Failed to retrieve product for update")
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}

	if existingProduct == nil {
		s.logger.Warn().
			Str("product_id", id.String()).
			Msg("Product not found for update")
		return nil, postgres.ErrResourceNotFound
	}

	// Check if name is being updated and if it conflicts with another product
	if dto.Name != nil && *dto.Name != existingProduct.Name {
		conflictingProduct, err := s.repo.GetByName(ctx, *dto.Name)
		if err != nil {
			s.logger.Error().Err(err).Msg("Error checking for conflicting product name")
			return nil, fmt.Errorf("error checking for conflicting product name: %w", err)
		}

		if conflictingProduct != nil && conflictingProduct.ID != id {
			s.logger.Warn().
				Str("product_id", id.String()).
				Str("conflicting_name", *dto.Name).
				Str("conflicting_id", conflictingProduct.ID.String()).
				Msg("Product name conflicts with existing product")
			return nil, fmt.Errorf("a product with the name '%s' already exists", *dto.Name)
		}
	}

	// Track if options are changing from empty to having values
	oldOptionsEmpty := len(existingProduct.Options) == 0
	var newOptionsEmpty bool
	var optionsChanged bool

	if dto.Options != nil {
		newOptionsEmpty = len(*dto.Options) == 0
		optionsChanged = true
	} else {
		newOptionsEmpty = oldOptionsEmpty
	}

	// Track if allow_subscription is changing to true
	var subscriptionEnabled bool
	if dto.AllowSubscription != nil {
		subscriptionEnabled = *dto.AllowSubscription
	} else {
		subscriptionEnabled = existingProduct.AllowSubscription
	}

	// Apply the updates to the existing product
	dto.ApplyToModel(existingProduct)

	// Update the product in the database
	err = s.repo.Update(ctx, existingProduct)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("product_id", id.String()).
			Msg("Failed to update product in database")
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Check if we need to trigger variant creation
	shouldCreateVariants := false
	
	// Case 1: Product now has options and allows subscription (and previously had no options)
	if oldOptionsEmpty && !newOptionsEmpty && subscriptionEnabled {
		shouldCreateVariants = true
		s.logger.Info().
			Str("product_id", id.String()).
			Msg("Product now has options and allows subscription - will trigger variant creation")
	}
	
	// Case 2: Product had no options, now has options (regardless of subscription)
	if oldOptionsEmpty && !newOptionsEmpty {
		shouldCreateVariants = true
		s.logger.Info().
			Str("product_id", id.String()).
			Msg("Product now has options - will trigger variant creation")
	}

	// Publish product updated event
	payload := events.ProductUpdatedPayload{
		ProductID:         existingProduct.ID.String(),
		Name:              existingProduct.Name,
		Description:       existingProduct.Description,
		ImageURL:          existingProduct.ImageURL,
		StockLevel:        existingProduct.StockLevel,
		Weight:            existingProduct.Weight,
		Origin:            existingProduct.Origin,
		RoastLevel:        existingProduct.RoastLevel,
		FlavorNotes:       existingProduct.FlavorNotes,
		Options:           existingProduct.Options,
		AllowSubscription: existingProduct.AllowSubscription,
		Active:            existingProduct.Active,
		Archived:          existingProduct.Archived,
		UpdatedAt:         existingProduct.UpdatedAt,
		OptionsChanged:    optionsChanged,
		ShouldCreateVariants: shouldCreateVariants,
	}

	err = s.eventBus.Publish(events.TopicProductUpdated, payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", id.String()).
			Msg("Failed to publish product updated event")
		// Don't return error since the product was updated successfully
	}

	s.logger.Info().
		Str("topic", events.TopicProductUpdated).
		Str("product_id", existingProduct.ID.String()).
		Bool("options_changed", optionsChanged).
		Bool("should_create_variants", shouldCreateVariants).
		Msg("Published product updated event")

	return existingProduct, nil
}

// Archive soft deletes a product by marking it as archived
func (s *productService) Archive(ctx context.Context, id uuid.UUID) error {
	s.logger.Info().
		Str("product_id", id.String()).
		Msg("Archiving product")

	// First, check if the product exists
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", id.String()).
			Msg("Error retrieving product for archiving")
		return fmt.Errorf("error retrieving product: %w", err)
	}

	if product == nil {
		s.logger.Warn().
			Str("product_id", id.String()).
			Msg("Product not found for archiving")
		return postgres.ErrResourceNotFound
	}

	// Check if product is already archived
	if product.Archived {
		s.logger.Info().
			Str("product_id", id.String()).
			Msg("Product is already archived")
		return nil // Already archived, nothing to do
	}

	// Archive the product
	err = s.repo.Archive(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", id.String()).
			Msg("Failed to archive product")
		return fmt.Errorf("failed to archive product: %w", err)
	}

	// Publish product archived event
	payload := map[string]interface{}{
		"product_id":  id.String(),
		"product_name": product.Name,
		"archived_at": time.Now().Format(time.RFC3339),
	}

	err = s.eventBus.Publish(events.TopicProductUpdated, payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", id.String()).
			Msg("Failed to publish product archived event")
		// Don't return the error since the product is already archived in DB
	}

	s.logger.Info().
		Str("product_id", id.String()).
		Msg("Product archived successfully")

	return nil
}

// Delete is now a dangerous operation that should only be used in specific cases
// such as removing test data. It attempts to hard delete a product from the database
func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Warn().
		Str("product_id", id.String()).
		Msg("Attempting hard delete of product - this should only be used for testing")

	// First, check if the product exists
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", id.String()).
			Msg("Error retrieving product for deletion")
		return fmt.Errorf("error retrieving product: %w", err)
	}

	if product == nil {
		s.logger.Warn().
			Str("product_id", id.String()).
			Msg("Product not found for deletion")
		return postgres.ErrResourceNotFound
	}

	// Delete the product from the database
	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", id.String()).
			Msg("Failed to delete product from database")
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// Publish product deleted event
	payload := map[string]string{
		"product_id": id.String(),
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	err = s.eventBus.Publish(events.TopicProductDeleted, payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", id.String()).
			Msg("Failed to publish product deleted event")
		// Don't return the error since the product is already deleted from DB
	}

	s.logger.Info().
		Str("product_id", id.String()).
		Msg("Product deleted successfully")

	return nil
}
