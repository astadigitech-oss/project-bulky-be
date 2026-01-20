-- =====================================================
-- ROLLBACK: Restore Dokumen Kebijakan Slug Trigger
-- =====================================================
-- Note: This rollback assumes slug column exists
-- If rolling back from migration 000085, restore slug column first
-- =====================================================

-- Recreate the function
CREATE OR REPLACE FUNCTION generate_slug_dokumen_kebijakan()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.slug IS NULL OR NEW.slug = '' THEN
        NEW.slug := LOWER(REGEXP_REPLACE(NEW.judul, '[^a-zA-Z0-9]+', '-', 'g'));
        NEW.slug := TRIM(BOTH '-' FROM NEW.slug);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Recreate the trigger
CREATE TRIGGER trg_generate_slug_dokumen
    BEFORE INSERT OR UPDATE ON dokumen_kebijakan
    FOR EACH ROW
    EXECUTE FUNCTION generate_slug_dokumen_kebijakan();
