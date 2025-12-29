-- migrations/000036_create_disclaimer.up.sql

CREATE TABLE disclaimer (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    judul VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE,
    konten TEXT NOT NULL,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_disclaimer_slug ON disclaimer(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_disclaimer_is_active ON disclaimer(is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_disclaimer_updated_at
    BEFORE UPDATE ON disclaimer
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Function to auto-generate slug
CREATE OR REPLACE FUNCTION generate_slug_disclaimer()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.slug IS NULL OR NEW.slug = '' THEN
        NEW.slug := LOWER(REGEXP_REPLACE(NEW.judul, '[^a-zA-Z0-9]+', '-', 'g'));
        NEW.slug := TRIM(BOTH '-' FROM NEW.slug);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_generate_slug_disclaimer
    BEFORE INSERT OR UPDATE ON disclaimer
    FOR EACH ROW
    EXECUTE FUNCTION generate_slug_disclaimer();

-- Trigger: Ensure single active disclaimer
CREATE OR REPLACE FUNCTION fn_ensure_single_active_disclaimer()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = true AND NEW.deleted_at IS NULL THEN
        UPDATE disclaimer 
        SET is_active = false, updated_at = NOW()
        WHERE id != NEW.id AND is_active = true AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_single_active_disclaimer
    AFTER INSERT OR UPDATE OF is_active ON disclaimer
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_disclaimer();

COMMENT ON TABLE disclaimer IS 'Halaman disclaimer. Hanya 1 yang bisa aktif.';
