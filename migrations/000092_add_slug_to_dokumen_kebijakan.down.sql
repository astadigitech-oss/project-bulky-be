-- Remove unique index
DROP INDEX IF EXISTS idx_dokumen_kebijakan_slug;

-- Remove slug column
ALTER TABLE dokumen_kebijakan
DROP COLUMN IF EXISTS slug;
