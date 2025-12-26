-- migrations/000033_create_informasi_pickup.up.sql

CREATE TABLE informasi_pickup (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    alamat TEXT NOT NULL,
    jam_operasional VARCHAR(100) NOT NULL,
    nomor_whatsapp VARCHAR(20) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    google_maps_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trg_informasi_pickup_updated_at
    BEFORE UPDATE ON informasi_pickup
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Singleton constraint: only allow 1 row
CREATE OR REPLACE FUNCTION check_informasi_pickup_singleton()
RETURNS TRIGGER AS $$
DECLARE
    v_count INT;
BEGIN
    SELECT COUNT(*) INTO v_count FROM informasi_pickup;
    IF v_count >= 1 THEN
        RAISE EXCEPTION 'Hanya boleh ada 1 informasi pickup (gunakan UPDATE, bukan INSERT)';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_informasi_pickup_singleton
    BEFORE INSERT ON informasi_pickup
    FOR EACH ROW
    EXECUTE FUNCTION check_informasi_pickup_singleton();

-- Seed default data
INSERT INTO informasi_pickup (alamat, jam_operasional, nomor_whatsapp) VALUES (
    'Jl. Cilodong Raya No.89, Cilodong, Kec. Cilodong, Kota Depok, Jawa Barat 16414',
    'Senin - Sabtu, 09.00 - 18.00 WIB',
    '62811833164'
);

COMMENT ON TABLE informasi_pickup IS 'Informasi lokasi pickup warehouse. Hanya 1 record (singleton).';
COMMENT ON COLUMN informasi_pickup.nomor_whatsapp IS 'Nomor WhatsApp, harus diawali 62 (tanpa +)';
