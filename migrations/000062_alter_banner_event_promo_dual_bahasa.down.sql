ALTER TABLE banner_event_promo DROP COLUMN IF EXISTS gambar_url_en;
ALTER TABLE banner_event_promo RENAME COLUMN gambar_url_id TO gambar;
