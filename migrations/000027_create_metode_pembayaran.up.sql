-- migrations/000027_create_metode_pembayaran.up.sql

CREATE TABLE metode_pembayaran (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL REFERENCES metode_pembayaran_group(id),
    nama VARCHAR(50) NOT NULL,
    kode VARCHAR(30) NOT NULL UNIQUE,
    logo VARCHAR(255),
    urutan INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_metode_pembayaran_group_id ON metode_pembayaran(group_id);
CREATE INDEX idx_metode_pembayaran_kode ON metode_pembayaran(kode);
CREATE INDEX idx_metode_pembayaran_urutan ON metode_pembayaran(urutan);

CREATE TRIGGER trg_metode_pembayaran_updated_at
    BEFORE UPDATE ON metode_pembayaran
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Seed data (Xendit payment methods)
INSERT INTO metode_pembayaran (group_id, nama, kode, urutan) VALUES
    -- Kartu Kredit
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'Kartu Kredit'), 'Kartu Kredit', 'CREDIT_CARD', 1),
    -- Bank Transfer / VA
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'BCA', 'BCA', 1),
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'Mandiri', 'MANDIRI', 2),
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'BNI', 'BNI', 3),
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'BRI', 'BRI', 4),
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'Permata', 'PERMATA', 5),
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'BSI', 'BSI', 6),
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'CIMB Niaga', 'CIMB', 7),
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'Bank Transfer / VA'), 'BJB', 'BJB', 8),
    -- QRIS
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'QRIS'), 'QRIS', 'QRIS', 1),
    -- E-Wallet
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'), 'GoPay', 'ID_GOPAY', 1),
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'), 'OVO', 'ID_OVO', 2),
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'), 'Dana', 'ID_DANA', 3),
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'E-Wallet'), 'LinkAja', 'ID_LINKAJA', 4),
    -- PayLater
    ((SELECT id FROM metode_pembayaran_group WHERE nama = 'PayLater'), 'Akulaku', 'ID_AKULAKU', 1);

COMMENT ON TABLE metode_pembayaran IS 'Metode pembayaran terintegrasi Xendit. Kode mengikuti Xendit channel code.';
