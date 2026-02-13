-- 000127_create_produk_merek.down.sql
-- Drop pivot table produk_merek

-- Drop indexes first
DROP INDEX IF EXISTS idx_produk_merek_merek_id;
DROP INDEX IF EXISTS idx_produk_merek_produk_id;

-- Drop the table
DROP TABLE IF EXISTS produk_merek;
