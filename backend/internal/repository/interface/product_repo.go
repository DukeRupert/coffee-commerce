// internal/repository/interfaces/product_repository.go
package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	GetByStripeID(ctx context.Context, stripeID string) (*model.Product, error)
	GetByName(ctx context.Context, name string) (*model.Product, error)
	// List retrieves all products, with optional filtering
	List(ctx context.Context, offset, limit int, includeInactive bool) ([]*model.Product, int, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStockLevel(ctx context.Context, id uuid.UUID, quantity int) error
}