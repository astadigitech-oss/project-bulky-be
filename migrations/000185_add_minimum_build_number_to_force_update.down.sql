ALTER TABLE force_update_app
DROP CONSTRAINT IF EXISTS chk_force_update_minimum_build_number,
DROP CONSTRAINT IF EXISTS chk_force_update_platform,
DROP COLUMN IF EXISTS minimum_build_number;

DROP TRIGGER IF EXISTS trg_single_active_force_update ON force_update_app;
CREATE TRIGGER trg_single_active_force_update
    AFTER INSERT OR UPDATE OF is_active ON force_update_app
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_force_update();
