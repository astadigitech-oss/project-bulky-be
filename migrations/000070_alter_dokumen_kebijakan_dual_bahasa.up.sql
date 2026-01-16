-- Add English language columns to dokumen_kebijakan table
ALTER TABLE dokumen_kebijakan
ADD COLUMN judul_en VARCHAR(200) NOT NULL DEFAULT '',
ADD COLUMN konten_en TEXT NOT NULL DEFAULT '';

-- Remove default values after adding columns
ALTER TABLE dokumen_kebijakan
ALTER COLUMN judul_en DROP DEFAULT,
ALTER COLUMN konten_en DROP DEFAULT;
