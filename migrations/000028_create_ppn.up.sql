-- migrations/000028_create_ppn.up.sql

CREATE TABLE ppn (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    persentase DECIMAL(5,2) NOT NULL,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TRIGGER trg_ppn_updated_at
    BEFORE UPDATE ON ppn
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger: Ensure single active PPN
CREATE OR REPLACE FUNCTION fn_ensure_single_active_ppn()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = true AND NEW.deleted_at IS NULL THEN
        UPDATE ppn 
        SET is_active = false, updated_at = NOW()
        WHERE id != NEW.id AND is_active = true AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_single_active_ppn
    AFTER INSERT OR UPDATE OF is_active ON ppn
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_ppn();

-- Seed default PPN
INSERT INTO ppn (persentase, is_active) VALUES (11.00, true);

COMMENT ON TABLE ppn IS 'Pengaturan PPN. Hanya 1 yang bisa aktif.';
