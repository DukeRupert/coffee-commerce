// internal/repository/postgres/price_repo.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// priceRepository implements the PriceRepository interface
type priceRepository struct {
	db     *DB
	logger zerolog.Logger
}

// NewPriceRepository creates a new PriceRepository
func NewPriceRepository(db *DB, logger *zerolog.Logger) *priceRepository {
	return &priceRepository{
		db:     db,
		logger: logger.With().Str("component", "price_repository").Logger(),
	}
}

// Create adds a new price to the database
func (r *priceRepository) Create(ctx context.Context, price *model.Price) error {
	query := `
        INSERT INTO prices (
            id, product_id, name, amount, currency, type,
            interval, interval_count, active, stripe_id,
            created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
        )
    `

	// Handle NULL values for interval and interval_count
	var interval interface{} = nil
	if price.Type == "recurring" && price.Interval != "" {
		interval = price.Interval
	}

	var intervalCount interface{} = nil
	if price.Type == "recurring" && price.IntervalCount > 0 {
		intervalCount = price.IntervalCount
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		price.ID,
		price.ProductID,
		price.Name,
		price.Amount,
		price.Currency,
		price.Type,
		interval,      // This will be NULL if interval is not set
		intervalCount, // This will be NULL if interval_count is not set
		price.Active,
		price.StripeID,
		price.CreatedAt,
		price.UpdatedAt,
	)

	if err != nil {
		r.logger.Error().Err(err).
			Str("price_id", price.ID.String()).
			Str("product_id", price.ProductID.String()).
			Msg("Failed to create price")
		return fmt.Errorf("failed to create price: %w", err)
	}

	return nil
}

// GetByID retrieves a price by its ID
func (r *priceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Price, error) {
	query := `
		SELECT
			id, product_id, name, amount, currency, type,
			interval, interval_count, active, stripe_id,
			created_at, updated_at
		FROM prices
		WHERE id = $1
	`

	var price model.Price
	var interval sql.NullString
	var intervalCount sql.NullInt32

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&price.ID,
		&price.ProductID,
		&price.Name,
		&price.Amount,
		&price.Currency,
		&price.Type,
		&interval,
		&intervalCount,
		&price.Active,
		&price.StripeID,
		&price.CreatedAt,
		&price.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Price not found
		}
		return nil, fmt.Errorf("failed to get price: %w", err)
	}

	// Handle nullable fields
	if interval.Valid {
		price.Interval = interval.String
	}

	if intervalCount.Valid {
		price.IntervalCount = int(intervalCount.Int32)
	}

	return &price, nil
}

// GetByProductID retrieves all prices for a product
func (r *priceRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]*model.Price, error) {
	query := `
		SELECT
			id, product_id, name, amount, currency, type,
			interval, interval_count, active, stripe_id,
			created_at, updated_at
		FROM prices
		WHERE product_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query prices: %w", err)
	}
	defer rows.Close()

	prices := make([]*model.Price, 0)

	for rows.Next() {
		var price model.Price
		var interval sql.NullString
		var intervalCount sql.NullInt32

		err := rows.Scan(
			&price.ID,
			&price.ProductID,
			&price.Name,
			&price.Amount,
			&price.Currency,
			&price.Type,
			&interval,
			&intervalCount,
			&price.Active,
			&price.StripeID,
			&price.CreatedAt,
			&price.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan price: %w", err)
		}

		// Handle nullable fields
		if interval.Valid {
			price.Interval = interval.String
		}

		if intervalCount.Valid {
			price.IntervalCount = int(intervalCount.Int32)
		}

		prices = append(prices, &price)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during price rows iteration: %w", err)
	}

	return prices, nil
}

// GetByStripeID retrieves a price by its Stripe ID
func (r *priceRepository) GetByStripeID(ctx context.Context, stripeID string) (*model.Price, error) {
	query := `
		SELECT
			id, product_id, name, amount, currency, type,
			interval, interval_count, active, stripe_id,
			created_at, updated_at
		FROM prices
		WHERE stripe_id = $1
	`

	var price model.Price
	var interval sql.NullString
	var intervalCount sql.NullInt32

	err := r.db.QueryRowContext(ctx, query, stripeID).Scan(
		&price.ID,
		&price.ProductID,
		&price.Name,
		&price.Amount,
		&price.Currency,
		&price.Type,
		&interval,
		&intervalCount,
		&price.Active,
		&price.StripeID,
		&price.CreatedAt,
		&price.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Price not found
		}
		return nil, fmt.Errorf("failed to get price by Stripe ID: %w", err)
	}

	// Handle nullable fields
	if interval.Valid {
		price.Interval = interval.String
	}

	if intervalCount.Valid {
		price.IntervalCount = int(intervalCount.Int32)
	}

	return &price, nil
}

// Update updates an existing price
func (r *priceRepository) Update(ctx context.Context, price *model.Price) error {
	price.UpdatedAt = time.Now()

	query := `
		UPDATE prices SET
			name = $1,
			amount = $2,
			currency = $3,
			type = $4,
			interval = $5,
			interval_count = $6,
			active = $7,
			stripe_id = $8,
			updated_at = $9
		WHERE id = $10
	`

	// Handle nullable fields
	var interval interface{} = nil
	if price.Interval != "" {
		interval = price.Interval
	}

	var intervalCount interface{} = nil
	if price.IntervalCount > 0 {
		intervalCount = price.IntervalCount
	}

	result, err := r.db.ExecContext(
		ctx,
		query,
		price.Name,
		price.Amount,
		price.Currency,
		price.Type,
		interval,
		intervalCount,
		price.Active,
		price.StripeID,
		price.UpdatedAt,
		price.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update price: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("price with ID %s not found", price.ID)
	}

	return nil
}

// Delete removes a price
func (r *priceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM prices WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete price: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("price with ID %s not found", id)
	}

	return nil
}
