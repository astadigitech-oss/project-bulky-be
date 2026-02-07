-- Rollback: Restore url_tujuan column
DROP INDEX IF EXISTS idx_banner_tujuan;

ALTER TABLE banner_event_promo 
ADD COLUMN url_tujuan VARCHAR(255) DEFAULT NULL;

ALTER TABLE banner_event_promo 
DROP COLUMN IF EXISTS tujuan;

COMMENT ON COLUMN banner_event_promo.url_tujuan IS 'Link redirect ketika banner diklik (nullable)';
