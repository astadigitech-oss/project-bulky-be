-- Rollback to previous version (from 000113)
DROP TRIGGER IF EXISTS trg_hero_section_auto_sync ON hero_section;
DROP FUNCTION IF EXISTS fn_hero_section_auto_sync();

CREATE OR REPLACE FUNCTION fn_hero_section_auto_sync()
RETURNS TRIGGER AS $$
BEGIN
    -- Rule 1: If tanggal_mulai AND tanggal_selesai are both set, auto-set is_default = true
    IF NEW.tanggal_mulai IS NOT NULL AND NEW.tanggal_selesai IS NOT NULL THEN
        NEW.is_default := true;
    END IF;
    
    -- Rule 2: If is_default is being set to true, clear tanggal from this record
    -- and unset is_default from other records
    IF NEW.is_default = true AND NEW.deleted_at IS NULL THEN
        -- Clear tanggal for this record (making it permanent default)
        NEW.tanggal_mulai := NULL;
        NEW.tanggal_selesai := NULL;
        
        -- Unset is_default from other records
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

CREATE TRIGGER trg_hero_section_auto_sync
    BEFORE INSERT OR UPDATE ON hero_section
    FOR EACH ROW
    EXECUTE FUNCTION fn_hero_section_auto_sync();
