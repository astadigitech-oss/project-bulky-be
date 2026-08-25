-- 000182_rename_reference_id_to_reference_code_in_produk.up.sql
ALTER TABLE produk RENAME COLUMN reference_id TO reference_code;
ALTER INDEX IF EXISTS idx_produk_reference_id RENAME TO idx_produk_reference_code;
COMMENT ON COLUMN produk.reference_code IS 'Reference Code for linking bundle products/cargos to WMS system data';
