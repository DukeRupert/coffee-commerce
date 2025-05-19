package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

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
