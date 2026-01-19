-- =====================================================
-- UPDATE: Dokumen Kebijakan Schema
-- =====================================================
-- Simplify to single language + add urutan field
-- Remove dual language support (judul_en, konten_en)
-- =====================================================

-- 1. Add urutan column
ALTER TABLE dokumen_kebijakan 
ADD COLUMN IF NOT EXISTS urutan INTEGER DEFAULT 0;

-- 2. Drop dual language columns if they exist
ALTER TABLE dokumen_kebijakan 
DROP COLUMN IF EXISTS judul_en,
DROP COLUMN IF EXISTS konten_en;

-- 3. Update slug to NOT NULL
ALTER TABLE dokumen_kebijakan 
ALTER COLUMN slug SET NOT NULL,
ALTER COLUMN slug TYPE VARCHAR(120);

-- 4. Update judul length
ALTER TABLE dokumen_kebijakan 
ALTER COLUMN judul TYPE VARCHAR(100);

-- 5. Update is_active default
ALTER TABLE dokumen_kebijakan 
ALTER COLUMN is_active SET DEFAULT true;

-- 6. Add comments
COMMENT ON COLUMN dokumen_kebijakan.konten IS 
'HTML content from rich text editor. Sanitized before storage.';

COMMENT ON COLUMN dokumen_kebijakan.urutan IS 
'Display order for fixed pages (1-7)';

COMMENT ON TABLE dokumen_kebijakan IS 
'Fixed policy pages (7 pages): Tentang Kami, Cara Membeli, Tentang Pembayaran, Hubungi Kami, FAQ, Syarat & Ketentuan, Kebijakan Privasi';

-- 7. Update existing data with urutan based on slug
UPDATE dokumen_kebijakan SET urutan = 1 WHERE slug = 'tentang-kami';
UPDATE dokumen_kebijakan SET urutan = 2 WHERE slug = 'cara-membeli';
UPDATE dokumen_kebijakan SET urutan = 3 WHERE slug = 'tentang-pembayaran';
UPDATE dokumen_kebijakan SET urutan = 4 WHERE slug = 'hubungi-kami';
UPDATE dokumen_kebijakan SET urutan = 5 WHERE slug = 'faq';
UPDATE dokumen_kebijakan SET urutan = 6 WHERE slug = 'syarat-ketentuan';
UPDATE dokumen_kebijakan SET urutan = 7 WHERE slug = 'kebijakan-privasi';
