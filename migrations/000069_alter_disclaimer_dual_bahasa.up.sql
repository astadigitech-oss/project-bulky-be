-- Add English language columns to disclaimer table
ALTER TABLE disclaimer
ADD COLUMN judul_en VARCHAR(200) NOT NULL DEFAULT '',
ADD COLUMN konten_en TEXT NOT NULL DEFAULT '';

-- Remove default values after adding columns
ALTER TABLE disclaimer
ALTER COLUMN judul_en DROP DEFAULT,
ALTER COLUMN konten_en DROP DEFAULT;
