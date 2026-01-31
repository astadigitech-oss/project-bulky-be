-- Rollback emergency fix - restore to old state
DO $$ 
BEGIN
    -- Drop new columns
    ALTER TABLE hero_section DROP COLUMN IF EXISTS tanggal_mulai;
    ALTER TABLE hero_section DROP COLUMN IF EXISTS tanggal_selesai;
    
    -- Rename is_default back to is_active if exists
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'hero_section' AND column_name = 'is_default'
    ) THEN
        ALTER TABLE hero_section RENAME COLUMN is_default TO is_active;
    END IF;
    
    -- Add urutan back if not exists
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'hero_section' AND column_name = 'urutan'
    ) THEN
        ALTER TABLE hero_section ADD COLUMN urutan INT DEFAULT 0;
    END IF;
END $$;

-- Drop new indexes and triggers
DROP INDEX IF EXISTS idx_hero_section_is_default;
DROP INDEX IF EXISTS idx_hero_section_tanggal;
DROP TRIGGER IF EXISTS trg_single_default_hero ON hero_section;
DROP FUNCTION IF EXISTS fn_ensure_single_default_hero();

-- Restore old trigger
CREATE OR REPLACE FUNCTION fn_ensure_single_active_hero()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = true AND NEW.deleted_at IS NULL THEN
        UPDATE hero_section 
        SET is_active = false,
            updated_at = NOW()
        WHERE id != NEW.id 
          AND is_active = true
          AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_single_active_hero ON hero_section;
CREATE TRIGGER trg_single_active_hero
    AFTER INSERT OR UPDATE OF is_active ON hero_section
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_hero();

-- Restore old indexes
DROP INDEX IF EXISTS idx_hero_section_is_active;
CREATE INDEX idx_hero_section_is_active ON hero_section(is_active) 
    WHERE deleted_at IS NULL;

DROP INDEX IF EXISTS idx_hero_section_urutan;
CREATE INDEX idx_hero_section_urutan ON hero_section(urutan) 
    WHERE is_active = true AND deleted_at IS NULL;
