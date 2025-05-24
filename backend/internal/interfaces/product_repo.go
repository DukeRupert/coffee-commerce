package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

type ProductRepository interface {
	// Core CRUD operations (currently implemented)
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	GetByName(ctx context.Context, name string) (*model.Product, error)
	GetByStripeID(ctx context.Context, stripeID string) (*model.Product, error)
	List(ctx context.Context, offset, limit int, includeInactive, includeArchived bool) ([]*model.Product, int, error)
	Update(ctx context.Context, product *model.Product) error
	Archive(ctx context.Context, id uuid.UUID) error // (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error  // (hard delete)
	UpdateStockLevel(ctx context.Context, id uuid.UUID, quantity int) error

	// Alternative lookup methods
	// GetBySKU(ctx context.Context, sku string) (*model.Product, error)
	// GetByBarcode(ctx context.Context, barcode string) (*model.Product, error)
	// GetBySlug(ctx context.Context, slug string) (*model.Product, error)

	// Batch operations
	// GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Product, error)
	// BulkCreate(ctx context.Context, products []*model.Product) error
	// BulkUpdate(ctx context.Context, products []*model.Product) error
	// BulkArchive(ctx context.Context, ids []uuid.UUID) error
	// BulkDelete(ctx context.Context, ids []uuid.UUID) error
	// BulkUpdateStockLevels(ctx context.Context, updates map[uuid.UUID]int) error

	// Search and filtering
	// Search(ctx context.Context, query string, offset, limit int) ([]*model.Product, int, error)
	// ListByOrigin(ctx context.Context, origin string, offset, limit int) ([]*model.Product, int, error)
	// ListByRoastLevel(ctx context.Context, roastLevel string, offset, limit int) ([]*model.Product, int, error)
	// ListByCategory(ctx context.Context, category string, offset, limit int) ([]*model.Product, int, error)
	// ListWithLowStock(ctx context.Context, threshold int) ([]*model.Product, error)
	// ListByPriceRange(ctx context.Context, minPrice, maxPrice int64, offset, limit int) ([]*model.Product, int, error)
	// ListByWeight(ctx context.Context, minWeight, maxWeight int, offset, limit int) ([]*model.Product, int, error)

	// Subscription-specific queries
	// ListSubscriptionProducts(ctx context.Context, offset, limit int) ([]*model.Product, int, error)
	// ListProductsWithOptions(ctx context.Context, offset, limit int) ([]*model.Product, int, error)
	// GetProductsNeedingVariants(ctx context.Context) ([]*model.Product, error)

	// Status and state queries
	// CountByStatus(ctx context.Context, active, archived bool) (int, error)
	// ListActive(ctx context.Context, offset, limit int) ([]*model.Product, int, error)
	// ListInactive(ctx context.Context, offset, limit int) ([]*model.Product, int, error)
	// ListArchived(ctx context.Context, offset, limit int) ([]*model.Product, int, error)
	// ListFeatured(ctx context.Context, limit int) ([]*model.Product, error)

	// Temporal queries
	// ListCreatedBetween(ctx context.Context, start, end time.Time, offset, limit int) ([]*model.Product, int, error)
	// ListUpdatedSince(ctx context.Context, since time.Time, offset, limit int) ([]*model.Product, int, error)
	// ListRecentlyArchived(ctx context.Context, since time.Time) ([]*model.Product, error)

	// Stock management queries
	// GetStockHistory(ctx context.Context, id uuid.UUID, from, to time.Time) ([]*model.StockAdjustment, error)
	// ListProductsWithStock(ctx context.Context, offset, limit int) ([]*model.Product, int, error)
	// ListOutOfStock(ctx context.Context) ([]*model.Product, error)
	// GetTotalInventoryValue(ctx context.Context) (int64, error)

	// Relationship queries
	// GetWithVariants(ctx context.Context, id uuid.UUID) (*model.Product, []*model.Variant, error)
	// GetWithPrices(ctx context.Context, id uuid.UUID) (*model.Product, []*model.Price, error)
	// GetRelatedProducts(ctx context.Context, id uuid.UUID) ([]*model.Product, error)

	// Analytics support
	// GetProductMetrics(ctx context.Context, id uuid.UUID, from, to time.Time) (*model.ProductMetrics, error)
	// ListTopSelling(ctx context.Context, limit int, period time.Duration) ([]*model.Product, error)
	// GetInventorySummary(ctx context.Context) (*model.InventorySummary, error)

	// Coffee-specific queries
	// ListByFarm(ctx context.Context, farmName string, offset, limit int) ([]*model.Product, int, error)
	// ListByCertification(ctx context.Context, certification string, offset, limit int) ([]*model.Product, int, error)
	// ListByProcessingMethod(ctx context.Context, method string, offset, limit int) ([]*model.Product, int, error)
	// GetByHarvestDate(ctx context.Context, from, to time.Time, offset, limit int) ([]*model.Product, int, error)

	// Validation and constraints
	// CheckNameExists(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error)
	// CheckSKUExists(ctx context.Context, sku string, excludeID *uuid.UUID) (bool, error)
	// ValidateProductConstraints(ctx context.Context, product *model.Product) error
}
