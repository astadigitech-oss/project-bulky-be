-- Remove English language columns from disclaimer table
ALTER TABLE disclaimer
DROP COLUMN IF EXISTS judul_en,
DROP COLUMN IF EXISTS konten_en;
