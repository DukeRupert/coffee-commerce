package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

type ProductRepository interface {
	// Create adds a new product to the database
	Create(ctx context.Context, product *model.Product) error

	// GetByID retrieves a product by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error)

	// GetByStripeID retrieves a product by its Stripe ID
	GetByStripeID(ctx context.Context, stripeID string) (*model.Product, error)

	// List retrieves all products, with optional filtering
	List(ctx context.Context, offset, limit int, includeInactive bool) ([]*model.Product, int, error)

	// Update updates an existing product
	Update(ctx context.Context, product *model.Product) error

	// Delete removes a product from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateStockLevel updates the stock level of a product
	UpdateStockLevel(ctx context.Context, id uuid.UUID, quantity int) error
}

// VariantRepository defines operations for managing product variants
type VariantRepository interface {
	// Create adds a new variant to the database
	Create(ctx context.Context, variant *model.Variant) error

	// GetByID retrieves a variant by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Variant, error)

	// GetByProductID retrieves all variants for a product
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*model.Variant, error)

	// Update updates an existing variant
	Update(ctx context.Context, variant *model.Variant) error

	// Delete removes a variant
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateStockLevel updates a variant's stock level
	UpdateStockLevel(ctx context.Context, id uuid.UUID, stockLevel int) error
}

// PriceRepository defines operations for managing price records
type PriceRepository interface {
	// Create adds a new price to the database
	Create(ctx context.Context, price *model.Price) error

	// GetByID retrieves a price by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Price, error)

	// GetByProductID retrieves all prices for a product
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)

	// GetByStripeID retrieves a price by its Stripe ID
	GetByStripeID(ctx context.Context, stripeID string) (*model.Price, error)

	// Update updates an existing price
	Update(ctx context.Context, price *model.Price) error

	// Delete removes a price
	Delete(ctx context.Context, id uuid.UUID) error
}
