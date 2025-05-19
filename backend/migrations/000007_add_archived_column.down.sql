
-- Remove archived column from products table
DROP INDEX IF EXISTS idx_products_archived;
ALTER TABLE products DROP COLUMN IF EXISTS archived;