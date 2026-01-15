-- Step 1: Rename nama to nama_id
ALTER TABLE kondisi_produk RENAME COLUMN nama TO nama_id;

-- Step 2: Add nama_en (nullable)
ALTER TABLE kondisi_produk ADD COLUMN nama_en VARCHAR(100) NULL;

-- Step 3: Add comment
COMMENT ON COLUMN kondisi_produk.nama_id IS 'Nama kondisi produk dalam Bahasa Indonesia';
COMMENT ON COLUMN kondisi_produk.nama_en IS 'Nama kondisi produk dalam Bahasa Inggris (opsional)';
