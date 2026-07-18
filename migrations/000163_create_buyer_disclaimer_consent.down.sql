-- migrations/000163_create_buyer_disclaimer_consent.down.sql

DROP INDEX IF EXISTS idx_bdc_disetujui_at;
DROP INDEX IF EXISTS idx_bdc_disclaimer_id;
DROP INDEX IF EXISTS idx_bdc_pesanan_id;
DROP INDEX IF EXISTS idx_bdc_buyer_id;
DROP TABLE IF EXISTS buyer_disclaimer_consent;
