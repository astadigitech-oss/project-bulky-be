-- Step 1: Rename gambar_url to gambar_url_id
ALTER TABLE banner_event_promo RENAME COLUMN gambar TO gambar_url_id;

-- Step 2: Add gambar_url_en (nullable)
ALTER TABLE banner_event_promo ADD COLUMN gambar_url_en VARCHAR(500) NULL;

-- Step 3: Add comment
COMMENT ON COLUMN banner_event_promo.gambar_url_id IS 'URL gambar banner Bahasa Indonesia';
COMMENT ON COLUMN banner_event_promo.gambar_url_en IS 'URL gambar banner Bahasa Inggris (opsional)';
