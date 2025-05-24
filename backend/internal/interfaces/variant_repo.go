package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

// VariantRepository defines operations for managing product variants
type VariantRepository interface {
	// Core CRUD operations (currently implemented)
	Create(ctx context.Context, variant *model.Variant) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Variant, error)
	GetByStripeID(ctx context.Context, stripeID string) (*model.Variant, error)
	GetByStripeProductID(ctx context.Context, stripeProductID string) (*model.Variant, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*model.Variant, error)
	Update(ctx context.Context, variant *model.Variant) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStockLevel(ctx context.Context, id uuid.UUID, stockLevel int) error

	// Batch operations
	// GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Variant, error)
	// BulkCreate(ctx context.Context, variants []*model.Variant) error
	// BulkUpdate(ctx context.Context, variants []*model.Variant) error
	// BulkDelete(ctx context.Context, ids []uuid.UUID) error
	// BulkUpdateStockLevels(ctx context.Context, updates map[uuid.UUID]int) error
	// BulkActivate(ctx context.Context, ids []uuid.UUID) error
	// BulkDeactivate(ctx context.Context, ids []uuid.UUID) error

	// Alternative lookup methods
	// GetByStripePriceID(ctx context.Context, stripePriceID string) (*model.Variant, error)
	// GetBySKU(ctx context.Context, sku string) (*model.Variant, error)
	// GetByPriceID(ctx context.Context, priceID uuid.UUID) ([]*model.Variant, error)

	// Option-based queries
	// GetByOptions(ctx context.Context, productID uuid.UUID, options map[string]string) (*model.Variant, error)
	// ListByOption(ctx context.Context, productID uuid.UUID, optionKey, optionValue string) ([]*model.Variant, error)
	// GetOptionValues(ctx context.Context, productID uuid.UUID, optionKey string) ([]string, error)
	// ListByWeight(ctx context.Context, weight int) ([]*model.Variant, error)

	// Status and availability queries
	// ListActive(ctx context.Context, productID uuid.UUID) ([]*model.Variant, error)
	// ListInactive(ctx context.Context, productID uuid.UUID) ([]*model.Variant, error)
	// ListAvailable(ctx context.Context, productID uuid.UUID) ([]*model.Variant, error) // active + in stock
	// ListOutOfStock(ctx context.Context, productID uuid.UUID) ([]*model.Variant, error)
	// ListWithLowStock(ctx context.Context, threshold int) ([]*model.Variant, error)

	// Stock management
	// ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error
	// ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) error
	// GetStockHistory(ctx context.Context, id uuid.UUID, from, to time.Time) ([]*model.StockAdjustment, error)
	// AdjustStock(ctx context.Context, id uuid.UUID, adjustment int, reason string) error

	// Subscription-specific queries
	// ListSubscriptionCompatible(ctx context.Context, productID uuid.UUID) ([]*model.Variant, error)
	// GetDefaultVariant(ctx context.Context, productID uuid.UUID) (*model.Variant, error)
	// ListByRecurringPrice(ctx context.Context, productID uuid.UUID) ([]*model.Variant, error)

	// Pricing queries
	// GetWithPrice(ctx context.Context, id uuid.UUID) (*model.Variant, *model.Price, error)
	// ListByPriceRange(ctx context.Context, productID uuid.UUID, minPrice, maxPrice int64) ([]*model.Variant, error)
	// ListWithoutPrices(ctx context.Context) ([]*model.Variant, error)

	// Validation and constraints
	// CheckOptionsUnique(ctx context.Context, productID uuid.UUID, options map[string]string, excludeID *uuid.UUID) (bool, error)
	// ValidateVariantConstraints(ctx context.Context, variant *model.Variant) error
	// CheckStripeIDExists(ctx context.Context, stripeProductID string, excludeID *uuid.UUID) (bool, error)

	// Analytics support
	// GetVariantMetrics(ctx context.Context, id uuid.UUID, from, to time.Time) (*model.VariantMetrics, error)
	// ListTopSelling(ctx context.Context, productID uuid.UUID, limit int, period time.Duration) ([]*model.Variant, error)
	// GetConversionRates(ctx context.Context, productID uuid.UUID) (map[uuid.UUID]float64, error)

	// Temporal queries
	// ListCreatedBetween(ctx context.Context, start, end time.Time) ([]*model.Variant, error)
	// ListUpdatedSince(ctx context.Context, since time.Time) ([]*model.Variant, error)
	// GetMostRecentVariant(ctx context.Context, productID uuid.UUID) (*model.Variant, error)
}
