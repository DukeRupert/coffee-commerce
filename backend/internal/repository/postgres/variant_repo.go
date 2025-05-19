// internal/repository/postgres/variant_repo.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// variantRepository implements the VariantRepository interface
type variantRepository struct {
	db     *DB
	logger zerolog.Logger
}

// NewVariantRepository creates a new VariantRepository
func NewVariantRepository(db *DB, logger *zerolog.Logger) *variantRepository {
	return &variantRepository{
		db:     db,
		logger: logger.With().Str("component", "variant_repository").Logger(),
	}
}

// internal/repository/postgres/variant_repo.go

// Create adds a new variant to the database
func (r *variantRepository) Create(ctx context.Context, variant *model.Variant) error {
	// Convert Options map to JSON string for storage
	optionsJSON, err := json.Marshal(variant.Options)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %w", err)
	}

	query := `
        INSERT INTO variants (
            id, product_id, price_id, stripe_price_id, weight,
            options, active, stock_level, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
        )
    `

	_, err = r.db.ExecContext(
		ctx,
		query,
		variant.ID,
		variant.ProductID,
		variant.PriceID,
		variant.StripePriceID,
		variant.Weight,
		optionsJSON,
		variant.Active,
		variant.StockLevel,
		variant.CreatedAt,
		variant.UpdatedAt,
	)

	if err != nil {
		r.logger.Error().Err(err).
			Str("variant_id", variant.ID.String()).
			Str("product_id", variant.ProductID.String()).
			Msg("Failed to create variant")
		return fmt.Errorf("failed to create variant: %w", err)
	}

	return nil
}

// GetByID retrieves a variant by its ID
func (r *variantRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Variant, error) {
	query := `
        SELECT
            id, product_id, price_id, stripe_price_id, weight,
            options, active, stock_level, created_at, updated_at
        FROM variants
        WHERE id = $1
    `

	var variant model.Variant
	var optionsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.PriceID,
		&variant.StripePriceID,
		&variant.Weight,
		&optionsJSON,
		&variant.Active,
		&variant.StockLevel,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Variant not found
		}
		return nil, fmt.Errorf("failed to get variant: %w", err)
	}

	// Unmarshal the options JSON
	if len(optionsJSON) > 0 {
		if err := json.Unmarshal(optionsJSON, &variant.Options); err != nil {
			return nil, fmt.Errorf("failed to unmarshal options: %w", err)
		}
	} else {
		// Initialize empty map if no options stored
		variant.Options = make(map[string]string)
	}

	return &variant, nil
}

// GetByProductID retrieves all variants for a product
func (r *variantRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]*model.Variant, error) {
	query := `
        SELECT
            id, product_id, price_id, stripe_price_id, weight,
            options, active, stock_level, created_at, updated_at
        FROM variants
        WHERE product_id = $1
        ORDER BY created_at
    `

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query variants: %w", err)
	}
	defer rows.Close()

	variants := make([]*model.Variant, 0)

	for rows.Next() {
		var variant model.Variant
		var optionsJSON []byte

		err := rows.Scan(
			&variant.ID,
			&variant.ProductID,
			&variant.PriceID,
			&variant.StripePriceID,
			&variant.Weight,
			&optionsJSON,
			&variant.Active,
			&variant.StockLevel,
			&variant.CreatedAt,
			&variant.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan variant: %w", err)
		}

		// Unmarshal the options JSON
		if len(optionsJSON) > 0 {
			if err := json.Unmarshal(optionsJSON, &variant.Options); err != nil {
				return nil, fmt.Errorf("failed to unmarshal options for variant %s: %w", variant.ID, err)
			}
		} else {
			// Initialize empty map if no options stored
			variant.Options = make(map[string]string)
		}

		variants = append(variants, &variant)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during variant rows iteration: %w", err)
	}

	return variants, nil
}

// Update updates an existing variant
func (r *variantRepository) Update(ctx context.Context, variant *model.Variant) error {
	variant.UpdatedAt = time.Now()

	// Convert Options map to JSON for storage
	optionsJSON, err := json.Marshal(variant.Options)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %w", err)
	}

	query := `
        UPDATE variants SET
            price_id = $1,
            stripe_price_id = $2,
            weight = $3,
            options = $5,
            active = $6,
            stock_level = $7,
            updated_at = $8
        WHERE id = $9
    `

	result, err := r.db.ExecContext(
		ctx,
		query,
		variant.PriceID,
		variant.StripePriceID,
		variant.Weight,
		optionsJSON,
		variant.Active,
		variant.StockLevel,
		variant.UpdatedAt,
		variant.ID,
	)

	if err != nil {
		r.logger.Error().Err(err).
			Str("variant_id", variant.ID.String()).
			Msg("Failed to update variant")
		return fmt.Errorf("failed to update variant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("variant with ID %s not found", variant.ID)
	}

	r.logger.Debug().
		Str("variant_id", variant.ID.String()).
		Int64("rows_affected", rowsAffected).
		Msg("Variant updated successfully")

	return nil
}

// Delete removes a variant
func (r *variantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM variants WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error().Err(err).
			Str("variant_id", id.String()).
			Msg("Failed to delete variant")
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("variant with ID %s not found", id)
	}

	r.logger.Debug().
		Str("variant_id", id.String()).
		Int64("rows_affected", rowsAffected).
		Msg("Variant deleted successfully")

	return nil
}

// UpdateStockLevel updates a variant's stock level
func (r *variantRepository) UpdateStockLevel(ctx context.Context, id uuid.UUID, stockLevel int) error {
	query := `
		UPDATE variants SET
			stock_level = $1,
			updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		stockLevel,
		time.Now(),
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update variant stock level: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("variant with ID %s not found", id)
	}

	return nil
}
