-- migrations/000035_create_dokumen_kebijakan.up.sql

CREATE TABLE dokumen_kebijakan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    judul VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE,
    konten TEXT NOT NULL,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_dokumen_kebijakan_slug ON dokumen_kebijakan(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_dokumen_kebijakan_is_active ON dokumen_kebijakan(is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_dokumen_kebijakan_updated_at
    BEFORE UPDATE ON dokumen_kebijakan
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Function to auto-generate slug from judul
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

CREATE TRIGGER trg_generate_slug_dokumen
    BEFORE INSERT OR UPDATE ON dokumen_kebijakan
    FOR EACH ROW
    EXECUTE FUNCTION generate_slug_dokumen_kebijakan();

COMMENT ON TABLE dokumen_kebijakan IS 'Dokumen kebijakan seperti Terms & Conditions, Privacy Policy, dll.';
