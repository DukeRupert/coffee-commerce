// internal/service/price_service.go
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dukerupert/coffee-commerce/internal/domain/dto"
	"github.com/dukerupert/coffee-commerce/internal/domain/model"
	"github.com/dukerupert/coffee-commerce/internal/events"
	"github.com/dukerupert/coffee-commerce/internal/interfaces"
	"github.com/dukerupert/coffee-commerce/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	stripeSDK "github.com/stripe/stripe-go/v82"
)

// priceService implements PriceService
type priceService struct {
	logger        zerolog.Logger
	eventBus      events.EventBus
	priceRepo     interfaces.PriceRepository
	productRepo   interfaces.ProductRepository
	variantRepo   interfaces.VariantRepository
	stripeService interfaces.StripeService
}

// NewPriceService creates a new price service
func NewPriceService(
	logger *zerolog.Logger,
	eventBus events.EventBus,
	priceRepo interfaces.PriceRepository,
	productRepo interfaces.ProductRepository,
	variantRepo interfaces.VariantRepository,
	stripeService interfaces.StripeService,
) interfaces.PriceService {
	subLogger := logger.With().Str("component", "price_service").Logger()
	return &priceService{
		logger:        subLogger,
		eventBus:      eventBus,
		priceRepo:     priceRepo,
		productRepo:   productRepo,
		variantRepo:   variantRepo,
		stripeService: stripeService,
	}
}

// Create creates a new price and optionally syncs it with Stripe
func (s *priceService) Create(ctx context.Context, createDTO *dto.PriceCreateDTO) (*model.Price, error) {
	s.logger.Info().
		Str("product_id", createDTO.ProductID.String()).
		Str("name", createDTO.Name).
		Int64("amount", createDTO.Amount).
		Str("currency", createDTO.Currency).
		Str("type", createDTO.Type).
		Msg("Creating price")

	// Verify that the product exists
	product, err := s.productRepo.GetByID(ctx, createDTO.ProductID)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", createDTO.ProductID.String()).
			Msg("Error retrieving product")
		return nil, fmt.Errorf("error retrieving product: %w", err)
	}

	if product == nil {
		s.logger.Warn().
			Str("product_id", createDTO.ProductID.String()).
			Msg("Product not found")
		return nil, postgres.ErrResourceNotFound
	}

	// Convert DTO to model
	price := createDTO.ToModel()

	// Create price in Stripe first (if Stripe is configured)
	var stripePrice *stripeSDK.Price
	if s.stripeService != nil {
		recurring := price.Type == "recurring"
		
		s.logger.Debug().
			Str("product_stripe_id", product.StripeID).
			Bool("recurring", recurring).
			Str("interval", price.Interval).
			Int("interval_count", price.IntervalCount).
			Msg("Creating price in Stripe")

		stripePrice, err = s.stripeService.CreatePrice(
			product.StripeID,
			price.Amount,
			price.Currency,
			recurring,
			price.Interval,
			int64(price.IntervalCount),
		)
		if err != nil {
			s.logger.Error().Err(err).
				Str("product_id", product.ID.String()).
				Str("product_stripe_id", product.StripeID).
				Msg("Failed to create price in Stripe")
			return nil, fmt.Errorf("failed to create price in Stripe: %w", err)
		}

		// Set the Stripe ID in our price model
		price.StripeID = stripePrice.ID
		
		s.logger.Info().
			Str("stripe_price_id", stripePrice.ID).
			Str("price_id", price.ID.String()).
			Msg("Successfully created price in Stripe")
	}

	// Save price to our database
	err = s.priceRepo.Create(ctx, price)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", price.ID.String()).
			Msg("Failed to save price to database")
		return nil, fmt.Errorf("failed to create price: %w", err)
	}

	// Publish price created event
	payload := map[string]interface{}{
		"price_id":       price.ID.String(),
		"product_id":     price.ProductID.String(),
		"name":           price.Name,
		"amount":         price.Amount,
		"currency":       price.Currency,
		"type":           price.Type,
		"interval":       price.Interval,
		"interval_count": price.IntervalCount,
		"active":         price.Active,
		"stripe_id":      price.StripeID,
		"created_at":     price.CreatedAt,
	}

	err = s.eventBus.Publish("prices.created", payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", price.ID.String()).
			Msg("Failed to publish price created event")
		// Don't return error since the price was already created
	}

	s.logger.Info().
		Str("price_id", price.ID.String()).
		Str("product_id", price.ProductID.String()).
		Msg("Price created successfully")

	return price, nil
}

// GetByID retrieves a price by its ID
func (s *priceService) GetByID(ctx context.Context, id uuid.UUID) (*model.Price, error) {
	s.logger.Debug().
		Str("price_id", id.String()).
		Msg("Retrieving price by ID")

	price, err := s.priceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", id.String()).
			Msg("Failed to retrieve price")
		return nil, fmt.Errorf("failed to retrieve price: %w", err)
	}

	if price == nil {
		s.logger.Warn().
			Str("price_id", id.String()).
			Msg("Price not found")
		return nil, postgres.ErrResourceNotFound
	}

	return price, nil
}

// GetByProductID retrieves all prices for a product
func (s *priceService) GetByProductID(ctx context.Context, productID uuid.UUID) ([]*model.Price, error) {
	s.logger.Debug().
		Str("product_id", productID.String()).
		Msg("Retrieving prices for product")

	// Verify that the product exists
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", productID.String()).
			Msg("Error retrieving product")
		return nil, fmt.Errorf("error retrieving product: %w", err)
	}

	if product == nil {
		s.logger.Warn().
			Str("product_id", productID.String()).
			Msg("Product not found")
		return nil, postgres.ErrResourceNotFound
	}

	prices, err := s.priceRepo.GetByProductID(ctx, productID)
	if err != nil {
		s.logger.Error().Err(err).
			Str("product_id", productID.String()).
			Msg("Failed to retrieve prices for product")
		return nil, fmt.Errorf("failed to retrieve prices: %w", err)
	}

	s.logger.Debug().
		Str("product_id", productID.String()).
		Int("price_count", len(prices)).
		Msg("Retrieved prices for product")

	return prices, nil
}

// Update updates an existing price
func (s *priceService) Update(ctx context.Context, id uuid.UUID, updateDTO *dto.PriceUpdateDTO) (*model.Price, error) {
	s.logger.Info().
		Str("price_id", id.String()).
		Msg("Updating price")

	// Get the existing price
	price, err := s.priceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", id.String()).
			Msg("Failed to retrieve price for update")
		return nil, fmt.Errorf("failed to retrieve price: %w", err)
	}

	if price == nil {
		s.logger.Warn().
			Str("price_id", id.String()).
			Msg("Price not found for update")
		return nil, postgres.ErrResourceNotFound
	}

	// Store original values for comparison
	originalAmount := price.Amount
	originalActive := price.Active

	// Apply the updates
	updateDTO.ApplyToModel(price)

	// Update in database
	err = s.priceRepo.Update(ctx, price)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", id.String()).
			Msg("Failed to update price in database")
		return nil, fmt.Errorf("failed to update price: %w", err)
	}

	// Publish price updated event
	payload := map[string]interface{}{
		"price_id":        price.ID.String(),
		"product_id":      price.ProductID.String(),
		"name":            price.Name,
		"amount":          price.Amount,
		"original_amount": originalAmount,
		"currency":        price.Currency,
		"type":            price.Type,
		"interval":        price.Interval,
		"interval_count":  price.IntervalCount,
		"active":          price.Active,
		"original_active": originalActive,
		"stripe_id":       price.StripeID,
		"updated_at":      price.UpdatedAt,
	}

	err = s.eventBus.Publish("prices.updated", payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", price.ID.String()).
			Msg("Failed to publish price updated event")
		// Don't return error since the price was already updated
	}

	s.logger.Info().
		Str("price_id", price.ID.String()).
		Msg("Price updated successfully")

	return price, nil
}

// Delete removes a price
func (s *priceService) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info().
		Str("price_id", id.String()).
		Msg("Deleting price")

	// Get the existing price for logging and validation
	price, err := s.priceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", id.String()).
			Msg("Failed to retrieve price for deletion")
		return fmt.Errorf("failed to retrieve price: %w", err)
	}

	if price == nil {
		s.logger.Warn().
			Str("price_id", id.String()).
			Msg("Price not found for deletion")
		return postgres.ErrResourceNotFound
	}

	// Check if any variants are using this price
	variants, err := s.variantRepo.GetByProductID(ctx, price.ProductID)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", id.String()).
			Str("product_id", price.ProductID.String()).
			Msg("Failed to check for variants using this price")
		return fmt.Errorf("failed to check for variants: %w", err)
	}

	// Count variants using this price
	variantsUsingPrice := 0
	for _, variant := range variants {
		if variant.PriceID == id {
			variantsUsingPrice++
		}
	}

	if variantsUsingPrice > 0 {
		s.logger.Warn().
			Str("price_id", id.String()).
			Int("variants_count", variantsUsingPrice).
			Msg("Cannot delete price - it is being used by variants")
		return fmt.Errorf("cannot delete price: %d variants are using this price", variantsUsingPrice)
	}

	// Delete the price
	err = s.priceRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", id.String()).
			Msg("Failed to delete price from database")
		return fmt.Errorf("failed to delete price: %w", err)
	}

	// Publish price deleted event
	payload := map[string]interface{}{
		"price_id":   id.String(),
		"product_id": price.ProductID.String(),
		"name":       price.Name,
		"stripe_id":  price.StripeID,
		"deleted_at": time.Now(),
	}

	err = s.eventBus.Publish("prices.deleted", payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", id.String()).
			Msg("Failed to publish price deleted event")
		// Don't return error since the price was already deleted
	}

	s.logger.Info().
		Str("price_id", id.String()).
		Msg("Price deleted successfully")

	return nil
}

// AssignToVariant assigns a price to a variant
func (s *priceService) AssignToVariant(ctx context.Context, assignmentDTO *dto.VariantPriceAssignmentDTO) error {
	s.logger.Info().
		Str("variant_id", assignmentDTO.VariantID.String()).
		Str("price_id", assignmentDTO.PriceID.String()).
		Msg("Assigning price to variant")

	// Verify that the price exists
	price, err := s.priceRepo.GetByID(ctx, assignmentDTO.PriceID)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", assignmentDTO.PriceID.String()).
			Msg("Error retrieving price")
		return fmt.Errorf("error retrieving price: %w", err)
	}

	if price == nil {
		s.logger.Warn().
			Str("price_id", assignmentDTO.PriceID.String()).
			Msg("Price not found")
		return postgres.ErrResourceNotFound
	}

	// Verify that the variant exists
	variant, err := s.variantRepo.GetByID(ctx, assignmentDTO.VariantID)
	if err != nil {
		s.logger.Error().Err(err).
			Str("variant_id", assignmentDTO.VariantID.String()).
			Msg("Error retrieving variant")
		return fmt.Errorf("error retrieving variant: %w", err)
	}

	if variant == nil {
		s.logger.Warn().
			Str("variant_id", assignmentDTO.VariantID.String()).
			Msg("Variant not found")
		return postgres.ErrResourceNotFound
	}

	// Verify that the price belongs to the variant's product
	if variant.ProductID != price.ProductID {
		s.logger.Warn().
			Str("variant_id", assignmentDTO.VariantID.String()).
			Str("variant_product_id", variant.ProductID.String()).
			Str("price_id", assignmentDTO.PriceID.String()).
			Str("price_product_id", price.ProductID.String()).
			Msg("Price and variant belong to different products")
		return fmt.Errorf("price belongs to a different product than the variant")
	}

	// Store the old price ID for the event
	oldPriceID := variant.PriceID

	// Update the variant with the new price
	variant.PriceID = assignmentDTO.PriceID
	variant.StripePriceID = price.StripeID
	variant.UpdatedAt = time.Now()

	err = s.variantRepo.Update(ctx, variant)
	if err != nil {
		s.logger.Error().Err(err).
			Str("variant_id", assignmentDTO.VariantID.String()).
			Str("price_id", assignmentDTO.PriceID.String()).
			Msg("Failed to update variant with new price")
		return fmt.Errorf("failed to assign price to variant: %w", err)
	}

	// Publish variant price assigned event
	payload := map[string]interface{}{
		"variant_id":     assignmentDTO.VariantID.String(),
		"product_id":     variant.ProductID.String(),
		"new_price_id":   assignmentDTO.PriceID.String(),
		"old_price_id":   oldPriceID.String(),
		"stripe_price_id": price.StripeID,
		"assigned_at":    time.Now(),
	}

	err = s.eventBus.Publish("variants.price_assigned", payload)
	if err != nil {
		s.logger.Error().Err(err).
			Str("variant_id", assignmentDTO.VariantID.String()).
			Str("price_id", assignmentDTO.PriceID.String()).
			Msg("Failed to publish variant price assigned event")
		// Don't return error since the assignment was successful
	}

	s.logger.Info().
		Str("variant_id", assignmentDTO.VariantID.String()).
		Str("price_id", assignmentDTO.PriceID.String()).
		Msg("Price assigned to variant successfully")

	return nil
}

// GetVariantsByPrice retrieves all variants using a specific price
func (s *priceService) GetVariantsByPrice(ctx context.Context, priceID uuid.UUID) ([]*model.Variant, error) {
	s.logger.Debug().
		Str("price_id", priceID.String()).
		Msg("Retrieving variants using price")

	// Verify that the price exists
	price, err := s.priceRepo.GetByID(ctx, priceID)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", priceID.String()).
			Msg("Error retrieving price")
		return nil, fmt.Errorf("error retrieving price: %w", err)
	}

	if price == nil {
		s.logger.Warn().
			Str("price_id", priceID.String()).
			Msg("Price not found")
		return nil, postgres.ErrResourceNotFound
	}

	// Get all variants for the product
	allVariants, err := s.variantRepo.GetByProductID(ctx, price.ProductID)
	if err != nil {
		s.logger.Error().Err(err).
			Str("price_id", priceID.String()).
			Str("product_id", price.ProductID.String()).
			Msg("Failed to retrieve variants for product")
		return nil, fmt.Errorf("failed to retrieve variants: %w", err)
	}

	// Filter variants that use this price
	var variantsUsingPrice []*model.Variant
	for _, variant := range allVariants {
		if variant.PriceID == priceID {
			variantsUsingPrice = append(variantsUsingPrice, variant)
		}
	}

	s.logger.Debug().
		Str("price_id", priceID.String()).
		Int("total_variants", len(allVariants)).
		Int("variants_using_price", len(variantsUsingPrice)).
		Msg("Retrieved variants using price")

	return variantsUsingPrice, nil
}

// ValidatePriceCompatibility checks if a price is compatible with a variant
func (s *priceService) ValidatePriceCompatibility(ctx context.Context, priceID, variantID uuid.UUID) error {
	// Get the price
	price, err := s.priceRepo.GetByID(ctx, priceID)
	if err != nil {
		return fmt.Errorf("error retrieving price: %w", err)
	}
	if price == nil {
		return postgres.ErrResourceNotFound
	}

	// Get the variant
	variant, err := s.variantRepo.GetByID(ctx, variantID)
	if err != nil {
		return fmt.Errorf("error retrieving variant: %w", err)
	}
	if variant == nil {
		return postgres.ErrResourceNotFound
	}

	// Check if they belong to the same product
	if price.ProductID != variant.ProductID {
		return fmt.Errorf("price and variant belong to different products")
	}

	// Check if the price is active
	if !price.Active {
		return fmt.Errorf("cannot assign inactive price to variant")
	}

	return nil
}

// SyncStripeProductIDs validates and fixes Stripe product ID mismatches using multiple search strategies
func (s *priceService) SyncStripeProductIDs(ctx context.Context) (*dto.SyncStripeProductIDsResult, error) {
	s.logger.Info().Msg("Starting intelligent Stripe product ID sync")

	result := &dto.SyncStripeProductIDsResult{
		Results: make([]dto.SyncResult, 0),
	}

	// Get all products from our database
	products, total, err := s.productRepo.List(ctx, 0, 100, true, true)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	result.TotalProducts = total

	// Get all Stripe products once for efficiency
	var allStripeProducts []*stripeSDK.Product
	if s.stripeService != nil {
		allStripeProducts, err = s.stripeService.ListAllProducts()
		if err != nil {
			return nil, fmt.Errorf("failed to list Stripe products: %w", err)
		}
		s.logger.Info().
			Int("stripe_products_count", len(allStripeProducts)).
			Msg("Retrieved all Stripe products for comparison")
	}

	for _, product := range products {
		syncResult := dto.SyncResult{
			ProductID:      product.ID.String(),
			ProductName:    product.Name,
			StoredStripeID: product.StripeID,
		}

		s.logger.Debug().
			Str("product_id", product.ID.String()).
			Str("product_name", product.Name).
			Str("stored_stripe_id", product.StripeID).
			Msg("Processing product for Stripe sync")

		if s.stripeService == nil {
			syncResult.Status = "error"
			syncResult.Error = "Stripe service not available"
			result.Summary.Errors++
			result.Results = append(result.Results, syncResult)
			continue
		}

		// Strategy 1: Try the stored ID first (might be correct)
		var foundStripeProduct *stripeSDK.Product
		var searchStrategy string

		if product.StripeID != "" {
			foundStripeProduct, err = s.stripeService.GetProduct(product.StripeID)
			if err == nil && foundStripeProduct != nil {
				searchStrategy = "stored_id"
				s.logger.Debug().
					Str("product_id", product.ID.String()).
					Str("stripe_id", foundStripeProduct.ID).
					Msg("Found Stripe product using stored ID")
			}
		}

		// Strategy 2: If stored ID didn't work, search by name
		if foundStripeProduct == nil {
			foundStripeProduct, err = s.stripeService.FindProductByName(product.Name)
			if err != nil {
				s.logger.Error().Err(err).
					Str("product_id", product.ID.String()).
					Str("product_name", product.Name).
					Msg("Error searching for Stripe product by name")
			} else if foundStripeProduct != nil {
				searchStrategy = "name_match"
				s.logger.Info().
					Str("product_id", product.ID.String()).
					Str("product_name", product.Name).
					Str("found_stripe_id", foundStripeProduct.ID).
					Msg("Found Stripe product by name match")
			}
		}

		// Strategy 3: Search by our product ID in Stripe metadata
		if foundStripeProduct == nil {
			foundStripeProduct, err = s.stripeService.FindProductByMetadata("original_product_id", product.ID.String())
			if err != nil {
				s.logger.Error().Err(err).
					Str("product_id", product.ID.String()).
					Msg("Error searching for Stripe product by metadata")
			} else if foundStripeProduct != nil {
				searchStrategy = "metadata_match"
				s.logger.Info().
					Str("product_id", product.ID.String()).
					Str("found_stripe_id", foundStripeProduct.ID).
					Msg("Found Stripe product by metadata match")
			}
		}

		// Strategy 4: Manual fuzzy matching by similar names
		if foundStripeProduct == nil {
			foundStripeProduct = s.findBestNameMatch(product.Name, allStripeProducts)
			if foundStripeProduct != nil {
				searchStrategy = "fuzzy_name_match"
				s.logger.Info().
					Str("product_id", product.ID.String()).
					Str("product_name", product.Name).
					Str("found_stripe_name", foundStripeProduct.Name).
					Str("found_stripe_id", foundStripeProduct.ID).
					Msg("Found Stripe product by fuzzy name match")
			}
		}

		// Process the results
		if foundStripeProduct == nil {
			syncResult.Status = "not_found"
			syncResult.Error = "No matching Stripe product found using any search strategy"
			result.Summary.NotFound++
			
			s.logger.Warn().
				Str("product_id", product.ID.String()).
				Str("product_name", product.Name).
				Str("stored_stripe_id", product.StripeID).
				Msg("Could not find matching Stripe product")
		} else {
			syncResult.ActualStripeID = foundStripeProduct.ID

			// Check if we need to update our database
			if foundStripeProduct.ID != product.StripeID {
				s.logger.Info().
					Str("product_id", product.ID.String()).
					Str("old_stripe_id", product.StripeID).
					Str("new_stripe_id", foundStripeProduct.ID).
					Str("search_strategy", searchStrategy).
					Msg("Stripe ID mismatch detected - updating database")

				syncResult.Status = "mismatch"
				result.Summary.Mismatches++

				// Update the product with correct Stripe ID
				product.StripeID = foundStripeProduct.ID
				updateErr := s.productRepo.Update(ctx, product)
				if updateErr != nil {
					s.logger.Error().Err(updateErr).
						Str("product_id", product.ID.String()).
						Msg("Failed to update product with correct Stripe ID")
					
					syncResult.Error = updateErr.Error()
					syncResult.Status = "error"
					result.Summary.Errors++
				} else {
					syncResult.Updated = true
					result.Summary.Updated++
					
					s.logger.Info().
						Str("product_id", product.ID.String()).
						Str("new_stripe_id", foundStripeProduct.ID).
						Str("strategy", searchStrategy).
						Msg("Successfully updated product with correct Stripe ID")
				}
			} else {
				syncResult.Status = "ok"
				result.Summary.OK++
				
				s.logger.Debug().
					Str("product_id", product.ID.String()).
					Str("stripe_id", foundStripeProduct.ID).
					Msg("Stripe ID is already correct")
			}
		}

		result.Results = append(result.Results, syncResult)
	}

	s.logger.Info().
		Interface("summary", result.Summary).
		Msg("Intelligent Stripe product ID sync completed")

	return result, nil
}

// findBestNameMatch performs fuzzy matching to find the best Stripe product match
func (s *priceService) findBestNameMatch(targetName string, stripeProducts []*stripeSDK.Product) *stripeSDK.Product {
	if len(stripeProducts) == 0 {
		return nil
	}

	targetName = strings.ToLower(strings.TrimSpace(targetName))
	
	// Look for partial matches, exact word matches, etc.
	var bestMatch *stripeSDK.Product
	bestScore := 0
	
	for _, stripeProduct := range stripeProducts {
		stripeName := strings.ToLower(strings.TrimSpace(stripeProduct.Name))
		
		// Score the match
		score := s.calculateNameMatchScore(targetName, stripeName)
		
		if score > bestScore && score >= 70 { // Minimum 70% confidence
			bestScore = score
			bestMatch = stripeProduct
		}
	}
	
	if bestMatch != nil {
		s.logger.Debug().
			Str("target_name", targetName).
			Str("matched_name", bestMatch.Name).
			Int("confidence_score", bestScore).
			Msg("Found fuzzy name match")
	}
	
	return bestMatch
}

// calculateNameMatchScore calculates a similarity score between two product names
func (s *priceService) calculateNameMatchScore(name1, name2 string) int {
	// Simple similarity scoring - you could use more sophisticated algorithms
	
	// Exact match
	if name1 == name2 {
		return 100
	}
	
	// Contains match
	if strings.Contains(name1, name2) || strings.Contains(name2, name1) {
		return 80
	}
	
	// Word overlap scoring
	words1 := strings.Fields(name1)
	words2 := strings.Fields(name2)
	
	if len(words1) == 0 || len(words2) == 0 {
		return 0
	}
	
	matches := 0
	for _, word1 := range words1 {
		for _, word2 := range words2 {
			if word1 == word2 {
				matches++
				break
			}
		}
	}
	
	// Score based on percentage of matching words
	maxWords := len(words1)
	if len(words2) > maxWords {
		maxWords = len(words2)
	}
	
	return (matches * 100) / maxWords
}