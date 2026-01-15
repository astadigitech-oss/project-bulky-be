-- Step 1: Rename gambar_url to gambar_url_id
ALTER TABLE hero_section RENAME COLUMN gambar TO gambar_url_id;

-- Step 2: Add gambar_url_en (nullable)
ALTER TABLE hero_section ADD COLUMN gambar_url_en VARCHAR(500) NULL;

-- Step 3: Add comment
COMMENT ON COLUMN hero_section.gambar_url_id IS 'URL gambar hero Bahasa Indonesia';
COMMENT ON COLUMN hero_section.gambar_url_en IS 'URL gambar hero Bahasa Inggris (opsional)';
