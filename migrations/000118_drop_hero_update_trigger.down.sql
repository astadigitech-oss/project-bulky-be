-- Rollback: Restore original trigger with INSERT OR UPDATE

DROP TRIGGER IF EXISTS trg_hero_section_auto_sync_insert ON hero_section;
DROP FUNCTION IF EXISTS fn_hero_section_auto_sync_insert();

-- Restore original function from 000115
CREATE OR REPLACE FUNCTION fn_hero_section_auto_sync()
RETURNS TRIGGER AS $$
BEGIN
    -- Rule 1: If tanggal_mulai AND tanggal_selesai are both set
    -- This hero becomes THE ONLY active hero (scheduled)
    IF NEW.tanggal_mulai IS NOT NULL AND NEW.tanggal_selesai IS NOT NULL THEN
        NEW.is_default := true;
        
        -- Clear ALL other records: unset is_default AND clear their dates
        UPDATE hero_section 
        SET is_default = false,
            tanggal_mulai = NULL,
            tanggal_selesai = NULL,
            updated_at = NOW()
        WHERE id != NEW.id 
          AND deleted_at IS NULL;
    
    -- Rule 2: If is_default = true (manual toggle to permanent default)
    -- This hero becomes THE ONLY active hero (permanent)
    ELSIF NEW.is_default = true AND NEW.deleted_at IS NULL THEN
        -- Clear tanggal for this record (permanent default)
        NEW.tanggal_mulai := NULL;
        NEW.tanggal_selesai := NULL;
        
        -- Clear ALL other records: unset is_default AND clear their dates
        UPDATE hero_section 
        SET is_default = false,
            tanggal_mulai = NULL,
            tanggal_selesai = NULL,
            updated_at = NOW()
        WHERE id != NEW.id 
          AND deleted_at IS NULL;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Recreate trigger with INSERT OR UPDATE
CREATE TRIGGER trg_hero_section_auto_sync
    BEFORE INSERT OR UPDATE ON hero_section
    FOR EACH ROW
    EXECUTE FUNCTION fn_hero_section_auto_sync();
