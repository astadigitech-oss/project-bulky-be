DROP TRIGGER IF EXISTS trg_kupon_updated_at ON kupon;
DROP INDEX IF EXISTS idx_kupon_kode_unique;
DROP INDEX IF EXISTS idx_kupon_is_active;
DROP INDEX IF EXISTS idx_kupon_tanggal_kedaluarsa;
DROP INDEX IF EXISTS idx_kupon_created_at;
DROP TABLE IF EXISTS kupon;
