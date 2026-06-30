-- Drop DB trigger that auto-logs status changes.
-- History is now handled exclusively by the application layer (Go code)
-- which provides richer context: changed_by and note fields.
DROP TRIGGER IF EXISTS trg_log_pesanan_status_change ON pesanan;
DROP FUNCTION IF EXISTS log_pesanan_status_change();
