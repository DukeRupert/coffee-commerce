-- Migration: 000008_create_sync_hashes_table.down.sql
-- Drop sync_hashes table and related indexes

DROP INDEX IF EXISTS idx_sync_hashes_variant_stripe_unique;
DROP INDEX IF EXISTS idx_sync_hashes_stripe_product_id;
DROP INDEX IF EXISTS idx_sync_hashes_variant_id;
DROP TABLE IF EXISTS sync_hashes;