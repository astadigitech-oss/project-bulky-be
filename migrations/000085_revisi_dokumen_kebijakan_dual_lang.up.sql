-- Migration: Add dual language support to dokumen_kebijakan
-- Remove slug (not needed for fixed 7 pages)

-- Add dual language columns
ALTER TABLE dokumen_kebijakan
ADD COLUMN judul_en VARCHAR(100),
ADD COLUMN konten_en TEXT;

-- Update existing data (set English = Indonesian temporarily)
UPDATE dokumen_kebijakan
SET judul_en = judul,
    konten_en = konten
WHERE judul_en IS NULL;

-- Set NOT NULL after data update
ALTER TABLE dokumen_kebijakan
ALTER COLUMN judul_en SET NOT NULL,
ALTER COLUMN konten_en SET NOT NULL;

-- Remove slug column (not needed for fixed pages)
ALTER TABLE dokumen_kebijakan
DROP COLUMN IF EXISTS slug;

-- Add comments
COMMENT ON COLUMN dokumen_kebijakan.judul IS 'Judul dalam Bahasa Indonesia';
COMMENT ON COLUMN dokumen_kebijakan.judul_en IS 'Judul dalam Bahasa Inggris';
COMMENT ON COLUMN dokumen_kebijakan.konten IS 'Konten HTML dalam Bahasa Indonesia';
COMMENT ON COLUMN dokumen_kebijakan.konten_en IS 'Konten HTML dalam Bahasa Inggris';
