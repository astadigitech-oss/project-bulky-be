-- migrations/000145_backfill_merek_produk_slug_en.up.sql
-- Backfill slug_en for merek_produk where slug_en is NULL but slug_id exists
-- All existing seed brands use the same name for ID and EN, so slug_en = slug_id

UPDATE merek_produk
SET slug_en = slug_id
WHERE slug_en IS NULL
  AND slug_id IS NOT NULL
  AND deleted_at IS NULL;
