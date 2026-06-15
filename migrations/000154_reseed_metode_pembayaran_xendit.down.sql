-- Rollback: hapus data hasil reseed
-- (data lama sebelum migration ini tidak dapat dikembalikan otomatis)
DELETE FROM metode_pembayaran;
DELETE FROM metode_pembayaran_group;
