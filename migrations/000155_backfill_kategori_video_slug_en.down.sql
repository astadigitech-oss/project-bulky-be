-- migrations/000155_backfill_kategori_video_slug_en.down.sql

UPDATE kategori_video SET slug_en = NULL WHERE deleted_at IS NULL;
