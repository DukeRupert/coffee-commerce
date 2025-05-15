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