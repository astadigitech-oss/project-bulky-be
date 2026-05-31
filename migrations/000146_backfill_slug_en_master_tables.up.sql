-- migrations/000146_backfill_slug_en_master_tables.up.sql
-- Backfill slug_en for sumber_produk, kondisi_paket, kondisi_produk

-- ============================================
-- 1. kondisi_paket
-- nama_id = nama_en for all records, so slug_en = slug_id
-- ============================================
UPDATE kondisi_paket
SET slug_en = slug_id
WHERE slug_en IS NULL
  AND slug_id IS NOT NULL
  AND deleted_at IS NULL;

-- ============================================
-- 2. kondisi_produk
-- nama_id = nama_en for all records, so slug_en = slug_id
-- ============================================
UPDATE kondisi_produk
SET slug_en = slug_id
WHERE slug_en IS NULL
  AND slug_id IS NOT NULL
  AND deleted_at IS NULL;

-- ============================================
-- 3. sumber_produk
-- nama_id != nama_en for some records, generate slug_en from nama_en
-- e.g. 'Local Supplier', 'Auction', 'Return', 'Liquidation', etc.
-- ============================================
UPDATE sumber_produk
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
