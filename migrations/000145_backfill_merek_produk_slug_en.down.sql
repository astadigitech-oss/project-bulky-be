-- migrations/000145_backfill_merek_produk_slug_en.down.sql
-- Rollback: clear slug_en for merek_produk (reset to NULL)
-- Note: this only clears records where slug_en = slug_id (i.e. auto-backfilled)

UPDATE merek_produk
SET slug_en = NULL
WHERE slug_en = slug_id
  AND deleted_at IS NULL;
