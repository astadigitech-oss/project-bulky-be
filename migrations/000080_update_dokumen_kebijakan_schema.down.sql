-- =====================================================
-- ROLLBACK: Dokumen Kebijakan Schema
-- =====================================================

-- 1. Add back dual language columns
ALTER TABLE dokumen_kebijakan 
ADD COLUMN IF NOT EXISTS judul_en VARCHAR(200),
ADD COLUMN IF NOT EXISTS konten_en TEXT;

-- 2. Remove urutan column
ALTER TABLE dokumen_kebijakan 
DROP COLUMN IF EXISTS urutan;

-- 3. Revert slug to nullable
ALTER TABLE dokumen_kebijakan 
ALTER COLUMN slug DROP NOT NULL,
ALTER COLUMN slug TYPE VARCHAR(200);

-- 4. Revert judul length
ALTER TABLE dokumen_kebijakan 
ALTER COLUMN judul TYPE VARCHAR(200);

-- 5. Revert is_active default
ALTER TABLE dokumen_kebijakan 
ALTER COLUMN is_active SET DEFAULT false;
