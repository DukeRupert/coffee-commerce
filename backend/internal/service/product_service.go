// internal/services/product_service.go
package services

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/dto"
	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/rs/zerolog"
)

// ProductService defines the interface for product business logic
type ProductService interface {
	Create(ctx context.Context, productDTO *dto.ProductCreateDTO) (*model.Product, error)
	List(ctx context.Context, page, pageSize int, includeInactive bool) ([]*model.Product, int, error)
}

// productService is the private implementation of ProductService
type productService struct {
	logger zerolog.Logger
}

// NewProductService creates a new instance of ProductService
func NewProductService(logger *zerolog.Logger) *productService {
	return &productService{
		logger: logger.With().Str("component", "product_service").Logger(),
	}
}

func (s *productService) Create(ctx context.Context, productDTO *dto.ProductCreateDTO) (*model.Product, error) {
	return nil, nil
}

func (s *productService) List(ctx context.Context, page, pageSize int, includeInactive bool) ([]*model.Product, int, error) {
	return nil, 0, nil
}
