package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/dto"
	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

// VariantService defines the interface for variant-related operations
type VariantService interface {
	// Methods will be added later
}

type ProductService interface {
	// Create(ctx context.Context, productDTO *dto.ProductCreateDTO) (*models.Product, error)
	Create(ctx context.Context, product *dto.ProductCreateDTO) (*model.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	List(ctx context.Context, offset, limit int, includeInactive, includeArchived bool) ([]*model.Product, int, error)
	Update(ctx context.Context, id uuid.UUID, productDTO *dto.ProductUpdateDTO) (*model.Product, error)
	Archive(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	// UpdateStockLevel(ctx context.Context, id uuid.UUID, quantity int) error
}
