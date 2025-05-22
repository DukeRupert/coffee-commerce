// backend/internal/sync/hash.go
package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
)

// StripeProductHashData represents the hashable data from a Stripe product
type StripeProductHashData struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Active      bool              `json:"active"`
	Images      []string          `json:"images"`
	Metadata    map[string]string `json:"metadata"`
}

// VariantHashData represents the hashable data from our variant
type VariantHashData struct {
	StripeProductID string            `json:"stripe_product_id"`
	StripePriceID   string            `json:"stripe_price_id"`
	Weight          int               `json:"weight"`
	Options         map[string]string `json:"options"`
	Active          bool              `json:"active"`
	StockLevel      int               `json:"stock_level"`
}

// ComputeStripeProductHash creates a consistent hash from Stripe product data
func ComputeStripeProductHash(stripeProduct stripe.Product) (string, error) {
	hashData := StripeProductHashData{
		ID:          stripeProduct.ID,
		Name:        stripeProduct.Name,
		Description: stripeProduct.Description,
		Active:      stripeProduct.Active,
		Images:      make([]string, len(stripeProduct.Images)),
		Metadata:    make(map[string]string),
	}

	// Copy and sort images for consistency
	copy(hashData.Images, stripeProduct.Images)
	sort.Strings(hashData.Images)

	// Copy metadata, excluding internal sync fields
	if stripeProduct.Metadata != nil {
		for k, v := range stripeProduct.Metadata {
			// Exclude our own sync-related metadata to prevent hash pollution
			if k != "sync_hash" && k != "last_sync" && k != "sync_source" {
				hashData.Metadata[k] = v
			}
		}
	}

	return computeHash(hashData)
}

// ComputeVariantHash creates a consistent hash from our variant data
func ComputeVariantHash(variant *model.Variant) (string, error) {
	hashData := VariantHashData{
		StripeProductID: variant.StripeProductID,
		StripePriceID:   variant.StripePriceID,
		Weight:          variant.Weight,
		Options:         make(map[string]string),
		Active:          variant.Active,
		StockLevel:      variant.StockLevel,
	}

	// Copy options map
	for k, v := range variant.Options {
		hashData.Options[k] = v
	}

	return computeHash(hashData)
}

// computeHash creates a SHA-256 hash of any data structure
func computeHash(data interface{}) (string, error) {
	// Convert to JSON for consistent serialization
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data for hashing: %w", err)
	}

	// Create hash
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:]), nil
}

// CreateSyncHashRecord creates a new sync hash record
func CreateSyncHashRecord(variantID uuid.UUID, stripeProductID, contentHash, syncSource string) *model.SyncHash {
	now := time.Now()
	return &model.SyncHash{
		ID:              uuid.New(),
		VariantID:       variantID,
		StripeProductID: stripeProductID,
		ContentHash:     contentHash,
		HashAlgorithm:   model.HashAlgorithmSHA256,
		SyncSource:      syncSource,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}