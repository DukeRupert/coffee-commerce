// backend/internal/repository/postgres/sync_repo.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// syncHashRepository implements the SyncHashRepository interface
type syncHashRepository struct {
	db     *DB
	logger zerolog.Logger
}

// NewSyncHashRepository creates a new SyncHashRepository
func NewSyncHashRepository(db *DB, logger *zerolog.Logger) *syncHashRepository {
	return &syncHashRepository{
		db:     db,
		logger: logger.With().Str("component", "sync_hash_repository").Logger(),
	}
}

// Create adds a new sync hash record
func (r *syncHashRepository) Create(ctx context.Context, syncHash *model.SyncHash) error {
	query := `
        INSERT INTO sync_hashes (
            id, variant_id, stripe_product_id, content_hash, 
            hash_algorithm, sync_source, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        )
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
		syncHash.ID,
		syncHash.VariantID,
		syncHash.StripeProductID,
		syncHash.ContentHash,
		syncHash.HashAlgorithm,
		syncHash.SyncSource,
		syncHash.CreatedAt,
		syncHash.UpdatedAt,
	)

	if err != nil {
		r.logger.Error().Err(err).
			Str("variant_id", syncHash.VariantID.String()).
			Str("stripe_product_id", syncHash.StripeProductID).
			Msg("Failed to create sync hash")
		return fmt.Errorf("failed to create sync hash: %w", err)
	}

	return nil
}

// GetByVariantID retrieves the latest sync hash for a variant
func (r *syncHashRepository) GetByVariantID(ctx context.Context, variantID uuid.UUID) (*model.SyncHash, error) {
	query := `
        SELECT id, variant_id, stripe_product_id, content_hash, 
               hash_algorithm, sync_source, created_at, updated_at
        FROM sync_hashes
        WHERE variant_id = $1
        ORDER BY updated_at DESC
        LIMIT 1
    `

	var syncHash model.SyncHash
	err := r.db.QueryRowContext(ctx, query, variantID).Scan(
		&syncHash.ID,
		&syncHash.VariantID,
		&syncHash.StripeProductID,
		&syncHash.ContentHash,
		&syncHash.HashAlgorithm,
		&syncHash.SyncSource,
		&syncHash.CreatedAt,
		&syncHash.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No sync hash found
		}
		return nil, fmt.Errorf("failed to get sync hash by variant ID: %w", err)
	}

	return &syncHash, nil
}

// GetByStripeProductID retrieves the latest sync hash for a Stripe product
func (r *syncHashRepository) GetByStripeProductID(ctx context.Context, stripeProductID string) (*model.SyncHash, error) {
	query := `
        SELECT id, variant_id, stripe_product_id, content_hash, 
               hash_algorithm, sync_source, created_at, updated_at
        FROM sync_hashes
        WHERE stripe_product_id = $1
        ORDER BY updated_at DESC
        LIMIT 1
    `

	var syncHash model.SyncHash
	err := r.db.QueryRowContext(ctx, query, stripeProductID).Scan(
		&syncHash.ID,
		&syncHash.VariantID,
		&syncHash.StripeProductID,
		&syncHash.ContentHash,
		&syncHash.HashAlgorithm,
		&syncHash.SyncSource,
		&syncHash.CreatedAt,
		&syncHash.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No sync hash found
		}
		return nil, fmt.Errorf("failed to get sync hash by Stripe product ID: %w", err)
	}

	return &syncHash, nil
}

// GetByVariantAndStripeID retrieves sync hash for specific variant-stripe product pair
func (r *syncHashRepository) GetByVariantAndStripeID(ctx context.Context, variantID uuid.UUID, stripeProductID string) (*model.SyncHash, error) {
	query := `
        SELECT id, variant_id, stripe_product_id, content_hash, 
               hash_algorithm, sync_source, created_at, updated_at
        FROM sync_hashes
        WHERE variant_id = $1 AND stripe_product_id = $2
        ORDER BY updated_at DESC
        LIMIT 1
    `

	var syncHash model.SyncHash
	err := r.db.QueryRowContext(ctx, query, variantID, stripeProductID).Scan(
		&syncHash.ID,
		&syncHash.VariantID,
		&syncHash.StripeProductID,
		&syncHash.ContentHash,
		&syncHash.HashAlgorithm,
		&syncHash.SyncSource,
		&syncHash.CreatedAt,
		&syncHash.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No sync hash found
		}
		return nil, fmt.Errorf("failed to get sync hash by variant and Stripe ID: %w", err)
	}

	return &syncHash, nil
}

// Upsert creates or updates a sync hash record
func (r *syncHashRepository) Upsert(ctx context.Context, syncHash *model.SyncHash) error {
	query := `
        INSERT INTO sync_hashes (
            id, variant_id, stripe_product_id, content_hash, 
            hash_algorithm, sync_source, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        )
        ON CONFLICT (variant_id, stripe_product_id) 
        DO UPDATE SET
            content_hash = EXCLUDED.content_hash,
            hash_algorithm = EXCLUDED.hash_algorithm,
            sync_source = EXCLUDED.sync_source,
            updated_at = EXCLUDED.updated_at
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
		syncHash.ID,
		syncHash.VariantID,
		syncHash.StripeProductID,
		syncHash.ContentHash,
		syncHash.HashAlgorithm,
		syncHash.SyncSource,
		syncHash.CreatedAt,
		syncHash.UpdatedAt,
	)

	if err != nil {
		r.logger.Error().Err(err).
			Str("variant_id", syncHash.VariantID.String()).
			Str("stripe_product_id", syncHash.StripeProductID).
			Msg("Failed to upsert sync hash")
		return fmt.Errorf("failed to upsert sync hash: %w", err)
	}

	return nil
}

// Delete removes a sync hash record
func (r *syncHashRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM sync_hashes WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sync hash: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sync hash with ID %s not found", id)
	}

	return nil
}

// DeleteByVariantID removes all sync hash records for a variant
func (r *syncHashRepository) DeleteByVariantID(ctx context.Context, variantID uuid.UUID) error {
	query := "DELETE FROM sync_hashes WHERE variant_id = $1"

	_, err := r.db.ExecContext(ctx, query, variantID)
	if err != nil {
		return fmt.Errorf("failed to delete sync hashes by variant ID: %w", err)
	}

	return nil
}

// GetHashHistory retrieves historical hashes for debugging/audit purposes
func (r *syncHashRepository) GetHashHistory(ctx context.Context, variantID uuid.UUID, limit int) ([]*model.SyncHash, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	query := `
        SELECT id, variant_id, stripe_product_id, content_hash, 
               hash_algorithm, sync_source, created_at, updated_at
        FROM sync_hashes
        WHERE variant_id = $1
        ORDER BY updated_at DESC
        LIMIT $2
    `

	rows, err := r.db.QueryContext(ctx, query, variantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query sync hash history: %w", err)
	}
	defer rows.Close()

	var hashes []*model.SyncHash
	for rows.Next() {
		var syncHash model.SyncHash
		err := rows.Scan(
			&syncHash.ID,
			&syncHash.VariantID,
			&syncHash.StripeProductID,
			&syncHash.ContentHash,
			&syncHash.HashAlgorithm,
			&syncHash.SyncSource,
			&syncHash.CreatedAt,
			&syncHash.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sync hash: %w", err)
		}
		hashes = append(hashes, &syncHash)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during sync hash history iteration: %w", err)
	}

	return hashes, nil
}