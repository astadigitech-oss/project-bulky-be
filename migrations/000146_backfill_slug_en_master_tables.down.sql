-- migrations/000146_backfill_slug_en_master_tables.down.sql
-- Rollback: clear slug_en for kondisi_paket, kondisi_produk, sumber_produk

UPDATE kondisi_paket
SET slug_en = NULL
WHERE slug_en = slug_id
  AND deleted_at IS NULL;

UPDATE kondisi_produk
SET slug_en = NULL
WHERE slug_en = slug_id
  AND deleted_at IS NULL;

UPDATE sumber_produk
SET slug_en = NULL
WHERE deleted_at IS NULL;
