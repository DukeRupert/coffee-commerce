package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

// PriceRepository defines operations for managing price records
type PriceRepository interface {
	// Core CRUD operations (currently implemented)
	Create(ctx context.Context, price *model.Price) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Price, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	GetByStripeID(ctx context.Context, stripeID string) (*model.Price, error)
	Update(ctx context.Context, price *model.Price) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Batch operations
	// GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Price, error)
	// BulkCreate(ctx context.Context, prices []*model.Price) error
	// BulkUpdate(ctx context.Context, prices []*model.Price) error
	// BulkDelete(ctx context.Context, ids []uuid.UUID) error
	// BulkActivate(ctx context.Context, ids []uuid.UUID) error
	// BulkDeactivate(ctx context.Context, ids []uuid.UUID) error

	// Price type queries
	// GetByType(ctx context.Context, productID uuid.UUID, priceType string) ([]*model.Price, error)
	// ListOneTimePrices(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// ListRecurringPrices(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// GetSubscriptionPrices(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)

	// Currency and amount queries
	// GetByCurrency(ctx context.Context, productID uuid.UUID, currency string) ([]*model.Price, error)
	// GetByAmountRange(ctx context.Context, productID uuid.UUID, minAmount, maxAmount int64) ([]*model.Price, error)
	// ListSupportedCurrencies(ctx context.Context) ([]string, error)
	// GetLowestPrice(ctx context.Context, productID uuid.UUID, currency string) (*model.Price, error)
	// GetHighestPrice(ctx context.Context, productID uuid.UUID, currency string) (*model.Price, error)

	// Subscription-specific queries
	// GetByInterval(ctx context.Context, productID uuid.UUID, interval string, intervalCount int) ([]*model.Price, error)
	// ListByInterval(ctx context.Context, interval string) ([]*model.Price, error)
	// GetMonthlyPrices(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// GetWeeklyPrices(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// GetAnnualPrices(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)

	// Status queries
	// ListActive(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// ListInactive(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// CountByStatus(ctx context.Context, productID uuid.UUID, active bool) (int, error)

	// Variant relationship queries
	// GetByVariantID(ctx context.Context, variantID uuid.UUID) (*model.Price, error)
	// ListUnassignedPrices(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// GetVariantsUsingPrice(ctx context.Context, priceID uuid.UUID) ([]*model.Variant, error)
	// CountVariantsUsingPrice(ctx context.Context, priceID uuid.UUID) (int, error)

	// Temporal queries
	// ListCreatedBetween(ctx context.Context, start, end time.Time) ([]*model.Price, error)
	// ListUpdatedSince(ctx context.Context, since time.Time) ([]*model.Price, error)
	// GetPriceHistory(ctx context.Context, priceID uuid.UUID, from, to time.Time) ([]*model.PriceHistory, error)

	// Validation and constraints
	// CheckDuplicatePrice(ctx context.Context, productID uuid.UUID, amount int64, currency, priceType string, excludeID *uuid.UUID) (bool, error)
	// ValidatePriceConstraints(ctx context.Context, price *model.Price) error
	// CheckStripeIDExists(ctx context.Context, stripeID string, excludeID *uuid.UUID) (bool, error)

	// Analytics support
	// GetPriceMetrics(ctx context.Context, priceID uuid.UUID, from, to time.Time) (*model.PriceMetrics, error)
	// GetAveragePrice(ctx context.Context, productID uuid.UUID, currency string) (int64, error)
	// GetPriceDistribution(ctx context.Context, productID uuid.UUID) (map[int64]int, error)

	// Promotional and dynamic pricing
	// ListPromotionalPrices(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// GetActivePromotions(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// ListExpiredPrices(ctx context.Context) ([]*model.Price, error)
	// GetPriceBeforeDate(ctx context.Context, priceID uuid.UUID, date time.Time) (*model.Price, error)

	// Regional pricing
	// GetByRegion(ctx context.Context, productID uuid.UUID, region string) ([]*model.Price, error)
	// ListRegions(ctx context.Context, productID uuid.UUID) ([]string, error)
}
