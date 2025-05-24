package interfaces

import (
	"context"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
)

// SyncHashRepository defines operations for managing sync hash records
type SyncHashRepository interface {
	// Core CRUD operations (currently implemented)
	Create(ctx context.Context, syncHash *model.SyncHash) error
	GetByVariantID(ctx context.Context, variantID uuid.UUID) (*model.SyncHash, error)
	GetByStripeProductID(ctx context.Context, stripeProductID string) (*model.SyncHash, error)
	GetByVariantAndStripeID(ctx context.Context, variantID uuid.UUID, stripeProductID string) (*model.SyncHash, error)
	Upsert(ctx context.Context, syncHash *model.SyncHash) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByVariantID(ctx context.Context, variantID uuid.UUID) error
	GetHashHistory(ctx context.Context, variantID uuid.UUID, limit int) ([]*model.SyncHash, error)

	// Batch operations
	// GetByVariantIDs(ctx context.Context, variantIDs []uuid.UUID) ([]*model.SyncHash, error)
	// GetByStripeProductIDs(ctx context.Context, stripeProductIDs []string) ([]*model.SyncHash, error)
	// BulkCreate(ctx context.Context, syncHashes []*model.SyncHash) error
	// BulkUpsert(ctx context.Context, syncHashes []*model.SyncHash) error
	// BulkDelete(ctx context.Context, ids []uuid.UUID) error

	// Hash comparison and validation
	// GetOutdatedHashes(ctx context.Context, algorithm string) ([]*model.SyncHash, error)
	// FindHashConflicts(ctx context.Context, variantID uuid.UUID) ([]*model.SyncHash, error)
	// ValidateHashIntegrity(ctx context.Context, variantID uuid.UUID) (bool, error)
	// GetHashesByAlgorithm(ctx context.Context, algorithm string) ([]*model.SyncHash, error)

	// Sync source tracking
	// GetBySyncSource(ctx context.Context, syncSource string) ([]*model.SyncHash, error)
	// GetWebhookSyncs(ctx context.Context, from, to time.Time) ([]*model.SyncHash, error)
	// GetAPISyncs(ctx context.Context, from, to time.Time) ([]*model.SyncHash, error)
	// CountBySyncSource(ctx context.Context, syncSource string) (int, error)

	// Temporal queries
	// GetRecentSyncs(ctx context.Context, since time.Time) ([]*model.SyncHash, error)
	// GetOldSyncs(ctx context.Context, before time.Time) ([]*model.SyncHash, error)
	// GetSyncActivity(ctx context.Context, from, to time.Time) ([]*model.SyncHash, error)
	// GetLastSyncTime(ctx context.Context, variantID uuid.UUID) (*time.Time, error)

	// Orphaned record management
	// FindOrphanedHashes(ctx context.Context) ([]*model.SyncHash, error)
	// CleanupOrphanedHashes(ctx context.Context) (int, error)
	// FindMissingSyncs(ctx context.Context) ([]uuid.UUID, error) // variant IDs without sync hashes

	// Analytics and monitoring
	// GetSyncStats(ctx context.Context, from, to time.Time) (*model.SyncStats, error)
	// GetFailedSyncs(ctx context.Context, from, to time.Time) ([]*model.SyncHash, error)
	// GetSyncFrequency(ctx context.Context, variantID uuid.UUID) (time.Duration, error)

	// Cleanup and maintenance
	// ArchiveOldHashes(ctx context.Context, before time.Time, keepLatest int) (int, error)
	// DeleteHistoryBefore(ctx context.Context, before time.Time) (int, error)
	// CompactHistory(ctx context.Context, variantID uuid.UUID, keepCount int) error

	// Conflict resolution
	// MarkHashConflict(ctx context.Context, syncHashID uuid.UUID, conflictReason string) error
	// ResolveHashConflict(ctx context.Context, syncHashID uuid.UUID, resolution string) error
	// ListUnresolvedConflicts(ctx context.Context) ([]*model.SyncHash, error)
}
