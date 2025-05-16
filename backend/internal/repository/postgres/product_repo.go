// internal/repositories/postgres/product_repository.go
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
	"github.com/rs/zerolog/log"
)

// ProductRepository implements the interfaces.ProductRepository interface
type productRepository struct {
	db     *DB
	logger zerolog.Logger
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *DB, logger *zerolog.Logger) *productRepository {
	return &productRepository{
		db:     db,
		logger: logger.With().Str("component", "product_repository").Logger(),
	}
}

// Create adds a new product to the database
func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	// Convert Options map to JSON string for storage
	optionsJSON, err := json.Marshal(product.Options)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %w", err)
	}

	query := `
		INSERT INTO products (
			id, name, description, image_url, active, stock_level,
			weight, origin, roast_level, flavor_notes, options, allow_subscription, stripe_id,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13,
			$14, $15
		)
	`

	// Access the underlying *sql.DB through our custom DB type
	_, err = r.db.ExecContext(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Description,
		product.ImageURL,
		product.Active,
		product.StockLevel,
		product.Weight,
		product.Origin,
		product.RoastLevel,
		product.FlavorNotes,
		optionsJSON,
		product.AllowSubscription,
		product.StripeID,
		product.CreatedAt,
		product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetByID retrieves a product by its ID
func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	query := `
		SELECT 
			id, name, description, image_url, active, stock_level,
			weight, origin, roast_level, flavor_notes, options, allow_subscription, stripe_id,
			created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product model.Product
	var optionsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.ImageURL,
		&product.Active,
		&product.StockLevel,
		&product.Weight,
		&product.Origin,
		&product.RoastLevel,
		&product.FlavorNotes,
		&optionsJSON,
		&product.AllowSubscription,
		&product.StripeID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Unmarshal the options JSON
	if len(optionsJSON) > 0 {
		if err := json.Unmarshal(optionsJSON, &product.Options); err != nil {
			return nil, fmt.Errorf("failed to unmarshal options: %w", err)
		}
	} else {
		// Initialize empty map if no options stored
		product.Options = make(map[string][]string)
	}

	return &product, nil
}

// GetProductByName retrieves a product by its name
func (r *productRepository) GetByName(ctx context.Context, name string) (*model.Product, error) {
    query := `
        SELECT id, stripe_id, name, description, image_url, origin, roast_level,
               stock_level, flavor_notes, active, options, allow_subscription,
               created_at, updated_at
        FROM products
        WHERE name = $1
        LIMIT 1
    `

    r.logger.Debug().Str("name", name).Msg("Querying product by name")

    var product model.Product
    var optionsJSON []byte // for storing the JSONB data

    err := r.db.QueryRowContext(ctx, query, name).Scan(
        &product.ID,
        &product.StripeID,
        &product.Name,
        &product.Description,
        &product.ImageURL,
        &product.Origin,
        &product.RoastLevel,
        &product.StockLevel,
        &product.FlavorNotes,
        &product.Active,
        &optionsJSON, // JSONB data needs special handling
        &product.AllowSubscription,
        &product.CreatedAt,
        &product.UpdatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            r.logger.Debug().Str("name", name).Msg("No product found with this name")
            return nil, nil // Return nil, nil to indicate no product found
        }
        r.logger.Error().Err(err).Str("name", name).Msg("Error querying product by name")
        return nil, fmt.Errorf("error querying product by name: %w", err)
    }

    // Parse the JSONB options data into the Options map
    if optionsJSON != nil && len(optionsJSON) > 0 {
        err = json.Unmarshal(optionsJSON, &product.Options)
        if err != nil {
            r.logger.Error().Err(err).Msg("Failed to unmarshal product options")
            return nil, fmt.Errorf("failed to unmarshal product options: %w", err)
        }
    }

    return &product, nil
}

// GetByStripeID retrieves a product by its Stripe ID
func (r *productRepository) GetByStripeID(ctx context.Context, stripeID string) (*model.Product, error) {
	query := `
		SELECT 
			id, name, description, image_url, active, stock_level,
			weight, origin, roast_level, flavor_notes, options, allow_subscription, stripe_id,
			created_at, updated_at
		FROM products
		WHERE stripe_id = $1
	`

	var product model.Product
	var optionsJSON []byte

	err := r.db.QueryRowContext(ctx, query, stripeID).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.ImageURL,
		&product.Active,
		&product.StockLevel,
		&product.Weight,
		&product.Origin,
		&product.RoastLevel,
		&product.FlavorNotes,
		&optionsJSON,
		&product.AllowSubscription,
		&product.StripeID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with Stripe ID %s not found", stripeID)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Unmarshal the options JSON
	if len(optionsJSON) > 0 {
		if err := json.Unmarshal(optionsJSON, &product.Options); err != nil {
			return nil, fmt.Errorf("failed to unmarshal options: %w", err)
		}
	} else {
		product.Options = make(map[string][]string)
	}

	return &product, nil
}

// List retrieves all products, with optional filtering
func (r *productRepository) List(ctx context.Context, offset, limit int, includeInactive bool) ([]*model.Product, int, error) {
	r.logger.Info().Msg("Executing List()")

	whereClause := ""
	if !includeInactive {
		whereClause = "WHERE active = true"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	listQuery := fmt.Sprintf(`
		SELECT 
			id, name, description, image_url, active, stock_level,
			weight, origin, roast_level, flavor_notes, options, allow_subscription, stripe_id,
			created_at, updated_at
		FROM products
		%s
		ORDER BY name
		LIMIT $1 OFFSET $2
	`, whereClause)

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	log.Debug().Int("total_count", total).Msg("Total record count")

	// If no products, return early
	if total == 0 {
		log.Debug().Msg("Total equals 0. Returning early")
		return []*model.Product{}, 0, nil
	}

	// Get products with pagination
	rows, err := r.db.QueryContext(ctx, listQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	products := make([]*model.Product, 0)
	for rows.Next() {
		var product model.Product
		var optionsJSON []byte

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.ImageURL,
			&product.Active,
			&product.StockLevel,
			&product.Weight,
			&product.Origin,
			&product.RoastLevel,
			&product.FlavorNotes,
			&optionsJSON,
			&product.AllowSubscription,
			&product.StripeID,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}

		// Unmarshal the options JSON
		if len(optionsJSON) > 0 {
			if err := json.Unmarshal(optionsJSON, &product.Options); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal options for product %s: %w", product.ID, err)
			}
		} else {
			product.Options = make(map[string][]string)
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during product rows iteration: %w", err)
	}

	return products, total, nil
}

// Update updates an existing product
func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	product.UpdatedAt = time.Now()

	// Convert Options map to JSON string for storage
	optionsJSON, err := json.Marshal(product.Options)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %w", err)
	}

	query := `
		UPDATE products SET
			name = $1,
			description = $2,
			image_url = $3,
			active = $4,
			stock_level = $5,
			weight = $6,
			origin = $7,
			roast_level = $8,
			flavor_notes = $9,
			options = $10,
			allow_subscription = $11,
			stripe_id = $12,
			updated_at = $13
		WHERE id = $14
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.ImageURL,
		product.Active,
		product.StockLevel,
		product.Weight,
		product.Origin,
		product.RoastLevel,
		product.FlavorNotes,
		optionsJSON,
		product.AllowSubscription,
		product.StripeID,
		product.UpdatedAt,
		product.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %s not found", product.ID)
	}

	return nil
}

// Delete removes a product from the database
func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM products WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %s not found", id)
	}

	return nil
}

// UpdateStockLevel updates the stock level of a product
func (r *productRepository) UpdateStockLevel(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `
		UPDATE products SET
			stock_level = $1,
			updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		quantity,
		time.Now(),
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update product stock level: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %s not found", id)
	}

	return nil
}