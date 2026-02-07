-- Rollback to original trigger
DROP TRIGGER IF EXISTS trg_hero_section_auto_sync ON hero_section;
DROP FUNCTION IF EXISTS fn_hero_section_auto_sync();

-- Restore original function
CREATE OR REPLACE FUNCTION fn_ensure_single_default_hero()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_default = true AND NEW.deleted_at IS NULL THEN
        UPDATE hero_section 
        SET is_default = false,
            updated_at = NOW()
        WHERE id != NEW.id 
          AND is_default = true
          AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_single_default_hero
    AFTER INSERT OR UPDATE OF is_default ON hero_section
    FOR EACH ROW
    WHEN (NEW.is_default = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_default_hero();
