package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/dto"
	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

// PriceService defines the interface for price-related operations
type PriceService interface {
	// Core CRUD operations (currently implemented)
	Create(ctx context.Context, createDTO *dto.PriceCreateDTO) (*model.Price, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Price, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	Update(ctx context.Context, id uuid.UUID, updateDTO *dto.PriceUpdateDTO) (*model.Price, error)
	Delete(ctx context.Context, id uuid.UUID) error
	AssignToVariant(ctx context.Context, assignmentDTO *dto.VariantPriceAssignmentDTO) error
	GetVariantsByPrice(ctx context.Context, priceID uuid.UUID) ([]*model.Variant, error)
	ValidatePriceCompatibility(ctx context.Context, priceID, variantID uuid.UUID) error
	SyncStripeProductIDs(ctx context.Context) (*dto.SyncStripeProductIDsResult, error)

	// Alternative lookup methods
	// GetByStripeID(ctx context.Context, stripeID string) (*model.Price, error)
	// GetActivePrices(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// GetPricesByType(ctx context.Context, productID uuid.UUID, priceType string) ([]*model.Price, error)

	// Price status management
	// Activate(ctx context.Context, id uuid.UUID) error
	// Deactivate(ctx context.Context, id uuid.UUID) error
	// BulkActivate(ctx context.Context, priceIDs []uuid.UUID) error
	// BulkDeactivate(ctx context.Context, priceIDs []uuid.UUID) error

	// Subscription pricing operations
	// CreateSubscriptionPrice(ctx context.Context, productID uuid.UUID, interval string, intervalCount int, amount int64, currency string) (*model.Price, error)
	// GetSubscriptionPrices(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// UpdateSubscriptionInterval(ctx context.Context, id uuid.UUID, interval string, intervalCount int) error
	// ConvertToSubscription(ctx context.Context, id uuid.UUID, interval string, intervalCount int) error
	// ConvertToOneTime(ctx context.Context, id uuid.UUID) error

	// Currency and localization
	// GetPricesByCurrency(ctx context.Context, productID uuid.UUID, currency string) ([]*model.Price, error)
	// CreateCurrencyVariant(ctx context.Context, basePriceID uuid.UUID, currency string, exchangeRate float64) (*model.Price, error)
	// UpdateExchangeRates(ctx context.Context, rates map[string]float64) error
	// GetSupportedCurrencies(ctx context.Context) ([]string, error)

	// Price tiers and volume discounts
	// CreateTieredPrice(ctx context.Context, productID uuid.UUID, tiers []dto.PriceTier) (*model.Price, error)
	// GetPriceForQuantity(ctx context.Context, priceID uuid.UUID, quantity int) (int64, error)
	// AddPriceTier(ctx context.Context, priceID uuid.UUID, tier dto.PriceTier) error
	// RemovePriceTier(ctx context.Context, priceID uuid.UUID, minQuantity int) error

	// Promotional pricing
	// CreatePromotionalPrice(ctx context.Context, basePriceID uuid.UUID, discountPercent float64, validFrom, validTo time.Time) (*model.Price, error)
	// ApplyDiscount(ctx context.Context, priceID uuid.UUID, discountAmount int64, reason string) error
	// RemoveDiscount(ctx context.Context, priceID uuid.UUID) error
	// GetActivePromotions(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)
	// SchedulePromotion(ctx context.Context, priceID uuid.UUID, startTime, endTime time.Time) error

	// Dynamic pricing
	// UpdateDynamicPrice(ctx context.Context, priceID uuid.UUID, newAmount int64, reason string) error
	// GetPriceHistory(ctx context.Context, priceID uuid.UUID, from, to time.Time) ([]*dto.PriceHistoryEntry, error)
	// RevertToPreviousPrice(ctx context.Context, priceID uuid.UUID) error
	// SetPriceRule(ctx context.Context, priceID uuid.UUID, rule dto.PriceRule) error

	// Bulk operations
	// BulkCreatePrices(ctx context.Context, prices []dto.PriceCreateDTO) ([]*model.Price, []error)
	// BulkUpdatePrices(ctx context.Context, updates map[uuid.UUID]dto.PriceUpdateDTO) error
	// BulkDeletePrices(ctx context.Context, priceIDs []uuid.UUID) error
	// ApplyBulkDiscount(ctx context.Context, productIDs []uuid.UUID, discountPercent float64) error

	// Stripe synchronization operations
	// SyncPriceWithStripe(ctx context.Context, priceID uuid.UUID) error
	// CreateStripePriceForPrice(ctx context.Context, priceID uuid.UUID) (string, error) // returns stripe price ID
	// UpdateStripePrice(ctx context.Context, priceID uuid.UUID) error
	// SyncAllPricesWithStripe(ctx context.Context) (*dto.StripePriceSyncResult, error)
	// ImportPricesFromStripe(ctx context.Context, productID uuid.UUID) ([]*model.Price, error)

	// Price comparison and analysis
	// ComparePrices(ctx context.Context, priceIDs []uuid.UUID) (*dto.PriceComparisonResult, error)
	// GetCompetitivePricing(ctx context.Context, productID uuid.UUID) (*dto.CompetitivePricingData, error)
	// AnalyzePricePerformance(ctx context.Context, priceID uuid.UUID, period time.Duration) (*dto.PricePerformanceMetrics, error)
	// GetPriceElasticity(ctx context.Context, productID uuid.UUID) (*dto.PriceElasticityData, error)

	// Validation and business rules
	// ValidatePriceRules(ctx context.Context, price *model.Price) error
	// CheckPriceConflicts(ctx context.Context, productID uuid.UUID, newPrice *dto.PriceCreateDTO) ([]string, error)
	// ValidateSubscriptionPricing(ctx context.Context, price *model.Price) error
	// EnforcePricingPolicies(ctx context.Context, priceID uuid.UUID) error

	// Regional and market-specific pricing
	// CreateRegionalPrice(ctx context.Context, basePriceID uuid.UUID, region string, adjustment float64) (*model.Price, error)
	// GetPriceForRegion(ctx context.Context, productID uuid.UUID, region string) (*model.Price, error)
	// UpdateRegionalAdjustment(ctx context.Context, priceID uuid.UUID, region string, adjustment float64) error
	// ListRegionalPrices(ctx context.Context, productID uuid.UUID) (map[string]*model.Price, error)

	// Tax and fee management
	// AddTaxToPrice(ctx context.Context, priceID uuid.UUID, taxRate float64, taxType string) error
	// RemoveTaxFromPrice(ctx context.Context, priceID uuid.UUID) error
	// CalculatePriceWithTax(ctx context.Context, priceID uuid.UUID, region string) (int64, error)
	// AddProcessingFee(ctx context.Context, priceID uuid.UUID, feeAmount int64) error

	// Coffee-specific pricing operations
	// CreateSeasonalPricing(ctx context.Context, productID uuid.UUID, season string, priceAdjustment float64) (*model.Price, error)
	// SetOriginPremium(ctx context.Context, productID uuid.UUID, origin string, premium int64) error
	// ApplyQualityGradeAdjustment(ctx context.Context, priceID uuid.UUID, grade string, adjustment float64) error
	// CreateFarmDirectPrice(ctx context.Context, productID uuid.UUID, farmID uuid.UUID, farmerShare float64) (*model.Price, error)

	// Import/Export operations
	// ImportPrices(ctx context.Context, data []dto.PriceImportDTO) (*dto.ImportResult, error)
	// ExportPrices(ctx context.Context, filters *dto.PriceExportFilters) ([]dto.PriceExportDTO, error)
	// ImportFromCSV(ctx context.Context, csvData []byte) (*dto.ImportResult, error)
	// ExportToCSV(ctx context.Context, filters *dto.PriceExportFilters) ([]byte, error)

	// Analytics and reporting
	// GetPricingMetrics(ctx context.Context, productID uuid.UUID, period time.Duration) (*dto.PricingMetrics, error)
	// GetRevenueForecast(ctx context.Context, priceID uuid.UUID, period time.Duration) (*dto.RevenueForecast, error)
	// GetPriceOptimizationSuggestions(ctx context.Context, productID uuid.UUID) ([]*dto.PriceOptimizationSuggestion, error)
	// GeneratePricingReport(ctx context.Context, filters *dto.PricingReportFilters) (*dto.PricingReport, error)

	// A/B testing for pricing
	// CreatePriceTest(ctx context.Context, productID uuid.UUID, testPrices []int64, testDuration time.Duration) (*dto.PriceTest, error)
	// GetPriceTestResults(ctx context.Context, testID uuid.UUID) (*dto.PriceTestResults, error)
	// EndPriceTest(ctx context.Context, testID uuid.UUID, winningPriceID uuid.UUID) error
	// ListActivePriceTests(ctx context.Context) ([]*dto.PriceTest, error)

	// Subscription lifecycle pricing
	// CreateTrialPrice(ctx context.Context, productID uuid.UUID, trialDuration time.Duration, trialAmount int64) (*model.Price, error)
	// CreateSetupFee(ctx context.Context, priceID uuid.UUID, setupAmount int64) error
	// ApplyLoyaltyDiscount(ctx context.Context, priceID uuid.UUID, customerID uuid.UUID, discountPercent float64) error
	// CalculateLifetimeValue(ctx context.Context, priceID uuid.UUID, churnRate float64) (int64, error)

	// Price protection and guarantees
	// CreatePriceProtection(ctx context.Context, priceID uuid.UUID, duration time.Duration) error
	// CheckPriceGuarantee(ctx context.Context, priceID uuid.UUID, competitorPrice int64) (bool, error)
	// ApplyPriceMatch(ctx context.Context, priceID uuid.UUID, matchPrice int64, source string) error
	// GetPriceProtectionStatus(ctx context.Context, priceID uuid.UUID) (*dto.PriceProtectionStatus, error)

	// Audit and compliance
	// GetPriceAuditLog(ctx context.Context, priceID uuid.UUID) ([]*dto.PriceAuditEntry, error)
	// ValidateComplianceRules(ctx context.Context, priceID uuid.UUID, jurisdiction string) error
	// GenerateComplianceReport(ctx context.Context, from, to time.Time) (*dto.ComplianceReport, error)
	// ArchivePriceChanges(ctx context.Context, before time.Time) error
}
