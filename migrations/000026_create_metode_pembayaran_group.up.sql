-- migrations/000026_create_metode_pembayaran_group.up.sql

CREATE TABLE metode_pembayaran_group (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(50) NOT NULL UNIQUE,
    urutan INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_metode_pembayaran_group_urutan ON metode_pembayaran_group(urutan);

CREATE TRIGGER trg_metode_pembayaran_group_updated_at
    BEFORE UPDATE ON metode_pembayaran_group
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Seed data
INSERT INTO metode_pembayaran_group (nama, urutan) VALUES
    ('Kartu Kredit', 1),
    ('Bank Transfer / VA', 2),
    ('QRIS', 3),
    ('E-Wallet', 4),
    ('PayLater', 5);

COMMENT ON TABLE metode_pembayaran_group IS 'Group metode pembayaran (VA, E-Wallet, QRIS, dll)';
