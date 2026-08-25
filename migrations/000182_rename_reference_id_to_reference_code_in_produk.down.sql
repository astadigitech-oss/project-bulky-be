-- 000182_rename_reference_id_to_reference_code_in_produk.down.sql
ALTER TABLE produk RENAME COLUMN reference_code TO reference_id;
ALTER INDEX IF EXISTS idx_produk_reference_code RENAME TO idx_produk_reference_id;
COMMENT ON COLUMN produk.reference_id IS 'Reference ID for linking bundle products to WMS system data';
