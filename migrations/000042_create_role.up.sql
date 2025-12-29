-- migrations/000042_create_role.up.sql

CREATE TABLE role (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(50) NOT NULL,
    kode VARCHAR(30) NOT NULL UNIQUE,
    deskripsi TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_role_kode ON role(kode);
CREATE INDEX idx_role_is_active ON role(is_active) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_role_updated_at
    BEFORE UPDATE ON role
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Seed default roles
INSERT INTO role (nama, kode, deskripsi) VALUES
    ('Super Admin', 'SUPER_ADMIN', 'Full access ke semua fitur sistem'),
    ('Admin', 'ADMIN', 'Akses ke sebagian besar fitur'),
    ('Staff', 'STAFF', 'Akses terbatas untuk operasional');

COMMENT ON TABLE role IS 'Role untuk admin users dengan permission berbeda';
