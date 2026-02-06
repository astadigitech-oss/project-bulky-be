-- Rollback: Remove dimensi & berat fields from produk table
DROP INDEX IF EXISTS idx_produk_berat;

ALTER TABLE produk
DROP COLUMN IF EXISTS panjang,
DROP COLUMN IF EXISTS lebar,
DROP COLUMN IF EXISTS tinggi,
DROP COLUMN IF EXISTS berat;
