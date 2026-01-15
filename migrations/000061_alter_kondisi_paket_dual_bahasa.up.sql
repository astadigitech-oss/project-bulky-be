-- Step 1: Rename nama to nama_id
ALTER TABLE kondisi_paket RENAME COLUMN nama TO nama_id;

-- Step 2: Add nama_en (nullable)
ALTER TABLE kondisi_paket ADD COLUMN nama_en VARCHAR(100) NULL;

-- Step 3: Add comment
COMMENT ON COLUMN kondisi_paket.nama_id IS 'Nama kondisi paket dalam Bahasa Indonesia';
COMMENT ON COLUMN kondisi_paket.nama_en IS 'Nama kondisi paket dalam Bahasa Inggris (opsional)';
