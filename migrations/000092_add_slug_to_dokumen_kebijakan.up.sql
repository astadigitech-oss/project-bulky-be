-- Add slug column to dokumen_kebijakan table
ALTER TABLE dokumen_kebijakan
ADD COLUMN slug VARCHAR(100);

-- Update existing records with slug based on urutan
UPDATE dokumen_kebijakan
SET slug = CASE urutan
    WHEN 1 THEN 'tentang-kami'
    WHEN 2 THEN 'cara-membeli'
    WHEN 3 THEN 'tentang-pembayaran'
    WHEN 4 THEN 'hubungi-kami'
    WHEN 5 THEN 'faq'
    WHEN 6 THEN 'syarat-ketentuan'
    WHEN 7 THEN 'kebijakan-privasi'
    ELSE 'dokumen-' || urutan::text
END;

-- Make slug NOT NULL and add unique constraint
ALTER TABLE dokumen_kebijakan
ALTER COLUMN slug SET NOT NULL;

CREATE UNIQUE INDEX idx_dokumen_kebijakan_slug ON dokumen_kebijakan(slug) WHERE deleted_at IS NULL;
