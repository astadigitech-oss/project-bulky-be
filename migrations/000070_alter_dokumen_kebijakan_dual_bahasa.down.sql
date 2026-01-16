-- Remove English language columns from dokumen_kebijakan table
ALTER TABLE dokumen_kebijakan
DROP COLUMN IF EXISTS judul_en,
DROP COLUMN IF EXISTS konten_en;
