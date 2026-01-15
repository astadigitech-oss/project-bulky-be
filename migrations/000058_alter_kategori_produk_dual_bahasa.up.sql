-- Step 1: Rename nama to nama_id
ALTER TABLE kategori_produk RENAME COLUMN nama TO nama_id;

-- Step 2: Add nama_en (nullable)
ALTER TABLE kategori_produk ADD COLUMN nama_en VARCHAR(100) NULL;

-- Step 3: Add comment
COMMENT ON COLUMN kategori_produk.nama_id IS 'Nama kategori dalam Bahasa Indonesia';
COMMENT ON COLUMN kategori_produk.nama_en IS 'Nama kategori dalam Bahasa Inggris (opsional)';
