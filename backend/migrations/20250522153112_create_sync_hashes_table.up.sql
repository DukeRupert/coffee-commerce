-- Migration: 000008_create_sync_hashes_table.up.sql
-- Create sync_hashes table for tracking content hashes to prevent unnecessary updates

CREATE TABLE sync_hashes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Entity references
    variant_id UUID REFERENCES variants(id) ON DELETE CASCADE,
    stripe_product_id VARCHAR(255) NOT NULL,
    
    -- Hash data
    content_hash VARCHAR(64) NOT NULL,
    hash_algorithm VARCHAR(20) DEFAULT 'sha256',
    
    -- Source tracking
    sync_source VARCHAR(50) NOT NULL, -- 'stripe_webhook', 'api_call', etc.
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for efficient lookups
CREATE INDEX idx_sync_hashes_variant_id ON sync_hashes(variant_id);
CREATE INDEX idx_sync_hashes_stripe_product_id ON sync_hashes(stripe_product_id);
CREATE UNIQUE INDEX idx_sync_hashes_variant_stripe_unique ON sync_hashes(variant_id, stripe_product_id);

-- Add comments for documentation
COMMENT ON TABLE sync_hashes IS 'Stores content hashes for sync state tracking between our variants and Stripe products';
COMMENT ON COLUMN sync_hashes.content_hash IS 'SHA-256 hash of the relevant data fields used for change detection';
COMMENT ON COLUMN sync_hashes.sync_source IS 'Source of the last sync operation that created this hash';
COMMENT ON COLUMN sync_hashes.hash_algorithm IS 'Algorithm used to generate the hash (future-proofing)';