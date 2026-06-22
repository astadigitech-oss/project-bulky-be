-- migrations/000155_backfill_kategori_video_slug_en.up.sql
-- Backfill slug_en for kategori_video from nama_en

UPDATE kategori_video
SET slug_en = regexp_replace(
    regexp_replace(
        lower(trim(nama_en)),
        '[^a-z0-9\-]+', '-', 'g'
    ),
    '-+', '-', 'g'
)
WHERE slug_en IS NULL
  AND nama_en IS NOT NULL
  AND nama_en != ''
  AND deleted_at IS NULL;
