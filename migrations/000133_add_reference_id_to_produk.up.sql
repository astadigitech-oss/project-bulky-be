-- 000133_add_reference_id_to_produk.up.sql
-- Add reference_id column to produk table for WMS system linking
ALTER TABLE produk ADD COLUMN reference_id VARCHAR(100) NULL;

-- Add index for reference_id for better query performance
CREATE INDEX idx_produk_reference_id ON produk(reference_id) WHERE reference_id IS NOT NULL;

-- Add comment for documentation
COMMENT ON COLUMN produk.reference_id IS 'Reference ID for linking bundle products to WMS system data';
