-- =====================================================
-- REVISI: banner_event_promo.tujuan
-- =====================================================
-- Mengubah url_tujuan (string) menjadi tujuan (JSONB)
-- Format: array of {id, slug} kategori produk
-- =====================================================

-- Step 1: Add new column tujuan (JSONB)
ALTER TABLE banner_event_promo 
ADD COLUMN tujuan JSONB DEFAULT NULL;

-- Step 2: Migrate existing data (if any url_tujuan contains kategori slug)
-- Skip jika url_tujuan berisi URL eksternal, set null
UPDATE banner_event_promo 
SET tujuan = NULL 
WHERE url_tujuan IS NOT NULL;

-- Step 3: Drop old column
ALTER TABLE banner_event_promo 
DROP COLUMN IF EXISTS url_tujuan;

-- Step 4: Add index for JSONB query (optional, for future filtering)
CREATE INDEX idx_banner_tujuan ON banner_event_promo 
USING GIN (tujuan) WHERE deleted_at IS NULL;

-- Step 5: Add comment
COMMENT ON COLUMN banner_event_promo.tujuan IS 'Array kategori produk tujuan [{id, slug}]. NULL = banner tanpa redirect.';
