DROP TRIGGER IF EXISTS trg_check_max_split_payment ON pesanan_pembayaran;
DROP TRIGGER IF EXISTS trg_pesanan_pembayaran_updated_at ON pesanan_pembayaran;
DROP FUNCTION IF EXISTS check_max_split_payment();
DROP TABLE IF EXISTS pesanan_pembayaran;
