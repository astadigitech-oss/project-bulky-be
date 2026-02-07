-- =====================================================
-- REVISI: banner_event_promo - Hapus is_active
-- =====================================================
-- Visibility sekarang hanya berdasarkan tanggal
-- =====================================================

-- Step 1: Drop related index
DROP INDEX IF EXISTS idx_banner_is_active;

-- Step 2: Drop column is_active
ALTER TABLE banner_event_promo 
DROP COLUMN IF EXISTS is_active;

-- Step 3: Update index untuk urutan (tanpa is_active condition)
DROP INDEX IF EXISTS idx_banner_urutan;
CREATE INDEX idx_banner_urutan ON banner_event_promo(urutan) 
    WHERE deleted_at IS NULL;

-- Step 4: Update comment
COMMENT ON TABLE banner_event_promo IS 'Banner event dan promo. Bisa multiple aktif (carousel). Visibility berdasarkan tanggal_mulai & tanggal_selesai.';
