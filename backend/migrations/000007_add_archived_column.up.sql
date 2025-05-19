-- Add archived column to products table
ALTER TABLE products ADD COLUMN archived BOOLEAN DEFAULT FALSE;

-- Create index for efficient querying of non-archived products
CREATE INDEX idx_products_archived ON products(archived);

-- Update existing deletion-related queries to filter by archived status
COMMENT ON TABLE products IS 'Products are never hard deleted, only archived via archived=true';