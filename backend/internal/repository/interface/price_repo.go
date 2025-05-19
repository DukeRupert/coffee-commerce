// internal/repository/interface/price_repo.go
package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

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
