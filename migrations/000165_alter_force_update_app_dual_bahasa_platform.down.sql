-- migrations/000165_alter_force_update_app_dual_bahasa_platform.down.sql

-- Kembalikan trigger single-active ke perilaku global (tanpa platform)
CREATE OR REPLACE FUNCTION fn_ensure_single_active_force_update()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = true AND NEW.deleted_at IS NULL THEN
        UPDATE force_update_app 
        SET is_active = false, updated_at = NOW()
        WHERE id != NEW.id AND is_active = true AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP INDEX IF EXISTS idx_force_update_platform_is_active;

ALTER TABLE force_update_app
DROP COLUMN IF EXISTS informasi_update_en,
DROP COLUMN IF EXISTS platform;

DROP TYPE IF EXISTS force_update_platform;
