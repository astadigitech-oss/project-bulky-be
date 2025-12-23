CREATE TABLE warehouse (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE,
    alamat TEXT,
    kota VARCHAR(100),
    kode_pos VARCHAR(10),
    telepon VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_warehouse_slug ON warehouse(slug);
CREATE INDEX idx_warehouse_is_active ON warehouse(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_warehouse_kota ON warehouse(kota);

CREATE TRIGGER update_warehouse_updated_at
    BEFORE UPDATE ON warehouse
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
