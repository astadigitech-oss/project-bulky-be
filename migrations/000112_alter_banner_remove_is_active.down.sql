-- Rollback: Restore is_active
ALTER TABLE banner_event_promo 
ADD COLUMN is_active BOOLEAN DEFAULT true;

-- Restore index
CREATE INDEX idx_banner_is_active ON banner_event_promo(is_active) 
    WHERE deleted_at IS NULL;

DROP INDEX IF EXISTS idx_banner_urutan;
CREATE INDEX idx_banner_urutan ON banner_event_promo(urutan) 
    WHERE is_active = true AND deleted_at IS NULL;
