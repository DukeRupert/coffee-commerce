package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/dto"
	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

type ProductService interface {
	// Core CRUD operations (currently implemented)
	Create(ctx context.Context, product *dto.ProductCreateDTO) (*model.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	List(ctx context.Context, offset, limit int, includeInactive, includeArchived bool) ([]*model.Product, int, error)
	Update(ctx context.Context, id uuid.UUID, productDTO *dto.ProductUpdateDTO) (*model.Product, error)
	Archive(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Alternative lookup methods
	// GetByName(ctx context.Context, name string) (*model.Product, error)
	// GetByStripeID(ctx context.Context, stripeID string) (*model.Product, error)
	// GetBySKU(ctx context.Context, sku string) (*model.Product, error)

	// Stock management operations
	// UpdateStockLevel(ctx context.Context, id uuid.UUID, quantity int) error
	// AdjustStock(ctx context.Context, id uuid.UUID, adjustment int, reason string) error
	// GetLowStockProducts(ctx context.Context, threshold int) ([]*model.Product, error)
	// BulkUpdateStockLevels(ctx context.Context, updates map[uuid.UUID]int) error
	// ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error
	// ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) error

	// Product status management
	// Activate(ctx context.Context, id uuid.UUID) error
	// Deactivate(ctx context.Context, id uuid.UUID) error
	// BulkActivate(ctx context.Context, productIDs []uuid.UUID) error
	// BulkDeactivate(ctx context.Context, productIDs []uuid.UUID) error
	// BulkArchive(ctx context.Context, productIDs []uuid.UUID) error

	// Product options and variants management
	// UpdateOptions(ctx context.Context, id uuid.UUID, options map[string][]string) error
	// AddOption(ctx context.Context, id uuid.UUID, optionKey string, optionValues []string) error
	// RemoveOption(ctx context.Context, id uuid.UUID, optionKey string) error
	// GetProductWithVariants(ctx context.Context, id uuid.UUID) (*dto.ProductWithVariantsDTO, error)
	// RegenerateVariants(ctx context.Context, id uuid.UUID) error

	// Search and filtering operations
	// Search(ctx context.Context, query string, filters *dto.ProductSearchFilters) ([]*model.Product, int, error)
	// ListByCategory(ctx context.Context, category string, offset, limit int) ([]*model.Product, int, error)
	// ListByOrigin(ctx context.Context, origin string, offset, limit int) ([]*model.Product, int, error)
	// ListByRoastLevel(ctx context.Context, roastLevel string, offset, limit int) ([]*model.Product, int, error)
	// ListFeatured(ctx context.Context, limit int) ([]*model.Product, error)
	// ListNewArrivals(ctx context.Context, limit int, since time.Time) ([]*model.Product, error)

	// Subscription-specific operations
	// ListSubscriptionProducts(ctx context.Context, offset, limit int) ([]*model.Product, int, error)
	// EnableSubscription(ctx context.Context, id uuid.UUID) error
	// DisableSubscription(ctx context.Context, id uuid.UUID) error
	// GetSubscriptionCompatibleProducts(ctx context.Context) ([]*model.Product, error)

	// Stripe synchronization operations
	// SyncWithStripe(ctx context.Context, id uuid.UUID) error
	// CreateStripeProduct(ctx context.Context, id uuid.UUID) (string, error) // returns stripe product ID
	// UpdateStripeProduct(ctx context.Context, id uuid.UUID) error
	// SyncAllProductsWithStripe(ctx context.Context) (*dto.StripeSyncResult, error)

	// Product duplication and templating
	// DuplicateProduct(ctx context.Context, sourceID uuid.UUID, newName string) (*model.Product, error)
	// CreateFromTemplate(ctx context.Context, templateID uuid.UUID, productData *dto.ProductCreateDTO) (*model.Product, error)

	// Image and media management
	// UpdateImages(ctx context.Context, id uuid.UUID, imageURLs []string) error
	// AddImage(ctx context.Context, id uuid.UUID, imageURL string) error
	// RemoveImage(ctx context.Context, id uuid.UUID, imageURL string) error
	// ReorderImages(ctx context.Context, id uuid.UUID, imageURLs []string) error

	// Pricing operations (if not handled by separate PriceService)
	// GetProductPrices(ctx context.Context, id uuid.UUID) ([]*model.Price, error)
	// SetBasePrice(ctx context.Context, id uuid.UUID, amount int64, currency string) error

	// Analytics and reporting
	// GetProductMetrics(ctx context.Context, id uuid.UUID, from, to time.Time) (*dto.ProductMetrics, error)
	// GetTopSellingProducts(ctx context.Context, limit int, period time.Duration) ([]*dto.ProductSales, error)
	// GetProductPerformance(ctx context.Context, productIDs []uuid.UUID, period time.Duration) ([]*dto.ProductPerformance, error)
	// GetInventoryReport(ctx context.Context) (*dto.InventoryReport, error)

	// Import/Export operations
	// ImportProducts(ctx context.Context, data []dto.ProductImportDTO) (*dto.ImportResult, error)
	// ExportProducts(ctx context.Context, filters *dto.ProductExportFilters) ([]dto.ProductExportDTO, error)
	// ImportFromCSV(ctx context.Context, csvData []byte) (*dto.ImportResult, error)
	// ExportToCSV(ctx context.Context, filters *dto.ProductExportFilters) ([]byte, error)

	// Validation and business rules
	// ValidateProductRules(ctx context.Context, product *model.Product) error
	// CheckNameAvailability(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error)
	// ValidateOptions(ctx context.Context, options map[string][]string) error
	// ValidateSubscriptionCompatibility(ctx context.Context, id uuid.UUID) error

	// Relationship management
	// AddRelatedProduct(ctx context.Context, productID, relatedProductID uuid.UUID, relationType string) error
	// RemoveRelatedProduct(ctx context.Context, productID, relatedProductID uuid.UUID) error
	// GetRelatedProducts(ctx context.Context, productID uuid.UUID) ([]*model.Product, error)
	// GetRecommendations(ctx context.Context, productID uuid.UUID, limit int) ([]*model.Product, error)

	// Seasonal and promotional operations
	// SetSeasonal(ctx context.Context, id uuid.UUID, seasonal bool, seasonStart, seasonEnd *time.Time) error
	// ListSeasonalProducts(ctx context.Context, season string) ([]*model.Product, error)
	// ScheduleActivation(ctx context.Context, id uuid.UUID, activateAt time.Time) error
	// ScheduleDeactivation(ctx context.Context, id uuid.UUID, deactivateAt time.Time) error

	// Coffee-specific operations
	// UpdateCuppingNotes(ctx context.Context, id uuid.UUID, cuppingNotes string) error
	// SetProcessingMethod(ctx context.Context, id uuid.UUID, processingMethod string) error
	// UpdateHarvestDate(ctx context.Context, id uuid.UUID, harvestDate time.Time) error
	// GetProductsByFarm(ctx context.Context, farmName string) ([]*model.Product, error)
	// GetProductsByCertification(ctx context.Context, certification string) ([]*model.Product, error)

	// Audit and history
	// GetProductHistory(ctx context.Context, id uuid.UUID) ([]*dto.ProductAuditLog, error)
	// GetStockHistory(ctx context.Context, id uuid.UUID, from, to time.Time) ([]*dto.StockAuditLog, error)
	// RestoreFromHistory(ctx context.Context, id uuid.UUID, version int) error
}
