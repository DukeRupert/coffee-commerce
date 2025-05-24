package interfaces

// import (
// 	"context"

// 	"github.com/dukerupert/coffee-commerce/internal/domain/dto"
// 	"github.com/dukerupert/coffee-commerce/internal/domain/model"
// 	"github.com/google/uuid"
// )

// VariantService defines the interface for variant-related operations
type VariantService interface {
	// Core CRUD operations that could be exposed via API

	// Create creates a new variant manually (alternative to event-driven creation)
	// Create(ctx context.Context, createDTO *dto.VariantCreateDTO) (*model.Variant, error)

	// GetByID retrieves a variant by its ID
	// GetByID(ctx context.Context, id uuid.UUID) (*model.Variant, error)

	// GetByProductID retrieves all variants for a specific product
	// GetByProductID(ctx context.Context, productID uuid.UUID) ([]*model.Variant, error)

	// GetByStripeProductID retrieves a variant by its Stripe product ID
	// GetByStripeProductID(ctx context.Context, stripeProductID string) (*model.Variant, error)

	// Update updates an existing variant
	// Update(ctx context.Context, id uuid.UUID, updateDTO *dto.VariantUpdateDTO) (*model.Variant, error)

	// Delete removes a variant
	// Delete(ctx context.Context, id uuid.UUID) error

	// Stock management operations

	// UpdateStockLevel updates the stock level for a variant
	// UpdateStockLevel(ctx context.Context, id uuid.UUID, stockLevel int) error

	// CheckStockAvailability checks if enough stock is available for an order
	// CheckStockAvailability(ctx context.Context, id uuid.UUID, quantity int) (bool, error)

	// ReserveStock reserves stock for a pending order
	// ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error

	// ReleaseStock releases previously reserved stock
	// ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) error

	// Variant generation and management

	// RegenerateVariants recreates all variants for a product based on current options
	// RegenerateVariants(ctx context.Context, productID uuid.UUID) error

	// CreateVariantsFromOptions creates variants for all combinations of the given options
	// CreateVariantsFromOptions(ctx context.Context, productID uuid.UUID, options map[string][]string) error

	// Pricing operations

	// AssignPrice assigns a price to a variant
	// AssignPrice(ctx context.Context, variantID, priceID uuid.UUID) error

	// CreatePriceForVariant creates a new price and assigns it to a variant
	// CreatePriceForVariant(ctx context.Context, variantID uuid.UUID, priceDTO *dto.PriceCreateDTO) (*model.Price, error)

	// Stripe synchronization operations

	// SyncWithStripe syncs a variant with its corresponding Stripe product and price
	// SyncWithStripe(ctx context.Context, variantID uuid.UUID) error

	// CreateStripeProductForVariant creates a Stripe product for a variant that doesn't have one
	// CreateStripeProductForVariant(ctx context.Context, variantID uuid.UUID) (*string, error) // returns stripe product ID

	// UpdateStripeProduct updates the Stripe product to match variant data
	// UpdateStripeProduct(ctx context.Context, variantID uuid.UUID) error

	// Bulk operations

	// BulkUpdateStockLevels updates stock levels for multiple variants
	// BulkUpdateStockLevels(ctx context.Context, updates map[uuid.UUID]int) error

	// BulkActivateVariants activates or deactivates multiple variants
	// BulkActivateVariants(ctx context.Context, variantIDs []uuid.UUID, active bool) error

	// Query and filtering operations

	// ListActiveVariants retrieves all active variants with pagination
	// ListActiveVariants(ctx context.Context, offset, limit int) ([]*model.Variant, int, error)

	// ListVariantsByOptions finds variants matching specific option criteria
	// ListVariantsByOptions(ctx context.Context, productID uuid.UUID, options map[string]string) ([]*model.Variant, error)

	// GetLowStockVariants retrieves variants with stock below the specified threshold
	// GetLowStockVariants(ctx context.Context, threshold int) ([]*model.Variant, error)

	// Validation operations

	// ValidateVariantOptions validates that variant options are valid for the product
	// ValidateVariantOptions(ctx context.Context, productID uuid.UUID, options map[string]string) error

	// CheckVariantUniqueness ensures a variant with the same options doesn't already exist
	// CheckVariantUniqueness(ctx context.Context, productID uuid.UUID, options map[string]string) (bool, error)

	// Reporting and analytics

	// GetVariantSalesData retrieves sales analytics for variants
	// GetVariantSalesData(ctx context.Context, variantID uuid.UUID, from, to time.Time) (*dto.VariantSalesData, error)

	// GetTopSellingVariants retrieves the best-performing variants
	// GetTopSellingVariants(ctx context.Context, limit int, period time.Duration) ([]*dto.VariantPerformance, error)

	// Subscription-specific operations

	// GetSubscriptionVariants retrieves variants that are available for subscription
	// GetSubscriptionVariants(ctx context.Context, productID uuid.UUID) ([]*model.Variant, error)

	// UpdateSubscriptionAvailability updates whether a variant can be used in subscriptions
	// UpdateSubscriptionAvailability(ctx context.Context, variantID uuid.UUID, available bool) error
}
