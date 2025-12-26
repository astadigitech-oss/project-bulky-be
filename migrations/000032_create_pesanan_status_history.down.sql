DROP TRIGGER IF EXISTS trg_log_pesanan_status_change ON pesanan;
DROP FUNCTION IF EXISTS log_pesanan_status_change();
DROP TABLE IF EXISTS pesanan_status_history;
DROP TYPE IF EXISTS status_history_type;
