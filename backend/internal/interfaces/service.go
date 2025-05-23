package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/dto"
	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
)

// VariantService defines the interface for variant-related operations
type VariantService interface {
	// Methods will be added later
}

type ProductService interface {
	Create(ctx context.Context, product *dto.ProductCreateDTO) (*model.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	List(ctx context.Context, offset, limit int, includeInactive, includeArchived bool) ([]*model.Product, int, error)
	Update(ctx context.Context, id uuid.UUID, productDTO *dto.ProductUpdateDTO) (*model.Product, error)
	Archive(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	// UpdateStockLevel(ctx context.Context, id uuid.UUID, quantity int) error
}

// PriceService defines the interface for price-related operations
type PriceService interface {
	// Create creates a new price and optionally syncs it with Stripe
	Create(ctx context.Context, createDTO *dto.PriceCreateDTO) (*model.Price, error)

	// GetByID retrieves a price by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Price, error)

	// GetByProductID retrieves all prices for a product
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)

	// Update updates an existing price
	Update(ctx context.Context, id uuid.UUID, updateDTO *dto.PriceUpdateDTO) (*model.Price, error)

	// Delete removes a price (only if not in use by variants)
	Delete(ctx context.Context, id uuid.UUID) error

	// AssignToVariant assigns a price to a variant
	AssignToVariant(ctx context.Context, assignmentDTO *dto.VariantPriceAssignmentDTO) error

	// GetVariantsByPrice retrieves all variants using a specific price
	GetVariantsByPrice(ctx context.Context, priceID uuid.UUID) ([]*model.Variant, error)

	// ValidatePriceCompatibility checks if a price is compatible with a variant
	ValidatePriceCompatibility(ctx context.Context, priceID, variantID uuid.UUID) error

	// SyncStripeProductIDs validates and fixes Stripe product ID mismatches
	SyncStripeProductIDs(ctx context.Context) (*dto.SyncStripeProductIDsResult, error)
}

type StripeService interface {
	CreateProduct(name, description string, imageURLs []string, metadata map[string]string) (*stripe.Product, error)
	CreatePrice(productID string, unitAmount int64, currency string, recurring bool,
		interval string, intervalCount int64) (*stripe.Price, error)
		GetProduct(productID string) (*stripe.Product, error)
	ListAllProducts() ([]*stripe.Product, error)
	FindProductByName(name string) (*stripe.Product, error)
	FindProductByMetadata(key, value string) (*stripe.Product, error)
}
