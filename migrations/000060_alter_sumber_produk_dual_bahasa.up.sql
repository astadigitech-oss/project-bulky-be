-- Step 1: Rename nama to nama_id
ALTER TABLE sumber_produk RENAME COLUMN nama TO nama_id;

-- Step 2: Add nama_en (nullable)
ALTER TABLE sumber_produk ADD COLUMN nama_en VARCHAR(100) NULL;

-- Step 3: Add comment
COMMENT ON COLUMN sumber_produk.nama_id IS 'Nama sumber produk dalam Bahasa Indonesia';
COMMENT ON COLUMN sumber_produk.nama_en IS 'Nama sumber produk dalam Bahasa Inggris (opsional)';
