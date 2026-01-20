-- Rollback: Restore slug and remove dual language

-- Add slug back
ALTER TABLE dokumen_kebijakan
ADD COLUMN slug VARCHAR(120);

-- Generate slugs from judul
UPDATE dokumen_kebijakan
SET slug = LOWER(REGEXP_REPLACE(judul, '[^a-zA-Z0-9]+', '-', 'g'))
WHERE slug IS NULL;

-- Set NOT NULL and unique constraint
ALTER TABLE dokumen_kebijakan
ALTER COLUMN slug SET NOT NULL;

ALTER TABLE dokumen_kebijakan
ADD CONSTRAINT dokumen_kebijakan_slug_key UNIQUE (slug);

-- Remove dual language columns
ALTER TABLE dokumen_kebijakan
DROP COLUMN IF EXISTS judul_en,
DROP COLUMN IF EXISTS konten_en;
