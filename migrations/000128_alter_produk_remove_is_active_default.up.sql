-- Remove default value from is_active column in produk table
-- This allows explicit control of is_active value from application code
ALTER TABLE produk ALTER COLUMN is_active DROP DEFAULT;
