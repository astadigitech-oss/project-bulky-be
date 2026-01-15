-- Step 1: Rename nama to nama_id
ALTER TABLE merek_produk RENAME COLUMN nama TO nama_id;

-- Step 2: Add nama_en (nullable)
ALTER TABLE merek_produk ADD COLUMN nama_en VARCHAR(100) NULL;

-- Step 3: Add comment
COMMENT ON COLUMN merek_produk.nama_id IS 'Nama merek dalam Bahasa Indonesia';
COMMENT ON COLUMN merek_produk.nama_en IS 'Nama merek dalam Bahasa Inggris (opsional)';
