-- 000133_add_reference_id_to_produk.down.sql
-- Drop index first
DROP INDEX IF EXISTS idx_produk_reference_id;

-- Drop reference_id column from produk table
ALTER TABLE produk DROP COLUMN IF EXISTS reference_id;
