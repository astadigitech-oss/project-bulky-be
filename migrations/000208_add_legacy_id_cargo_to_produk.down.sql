DROP INDEX IF EXISTS idx_produk_legacy_id_cargo;

ALTER TABLE produk
    DROP COLUMN IF EXISTS legacy_id_cargo;
