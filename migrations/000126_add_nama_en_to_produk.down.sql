DROP INDEX IF EXISTS idx_produk_nama_en;
ALTER TABLE produk DROP COLUMN IF EXISTS nama_en;

-- Rename "nama_id" back to "nama"
ALTER TABLE produk
RENAME COLUMN nama_id TO nama;

COMMENT ON COLUMN produk.nama IS 'Nama produk';