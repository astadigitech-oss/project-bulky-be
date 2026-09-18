-- ============================================================
-- PRD 02 — Rollback tabel auction
-- CATATAN: down migration ini menghapus data lelang.
-- Prosedur rollback production default adalah rollback aplikasi +
-- disable fitur, BUKAN menjalankan down migration ini (PRD 02).
-- ============================================================

DROP TABLE IF EXISTS auction_idempotency;
DROP TABLE IF EXISTS auction_audit_logs;
DROP TABLE IF EXISTS auction_events;
DROP TABLE IF EXISTS auction_stock_reservations;
DROP TABLE IF EXISTS auction_winners;
DROP TABLE IF EXISTS auction_bids;
DROP TABLE IF EXISTS auction_shipping_quotes;
DROP TABLE IF EXISTS auction_batch_assets;
DROP TABLE IF EXISTS auction_assets;
DROP TABLE IF EXISTS auction_batch_brands;
DROP TABLE IF EXISTS auction_batch_items;
DROP TABLE IF EXISTS auction_batches;
