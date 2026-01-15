ALTER TABLE hero_section DROP COLUMN IF EXISTS gambar_url_en;
ALTER TABLE hero_section RENAME COLUMN gambar_url_id TO gambar;
