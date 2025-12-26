-- migrations/000034_create_jadwal_gudang.up.sql

CREATE TABLE jadwal_gudang (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    informasi_pickup_id UUID NOT NULL REFERENCES informasi_pickup(id) ON DELETE CASCADE,
    hari INT NOT NULL CHECK (hari >= 0 AND hari <= 6),
    jam_buka TIME,
    jam_tutup TIME,
    is_buka BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(informasi_pickup_id, hari)
);

CREATE INDEX idx_jadwal_gudang_pickup_id ON jadwal_gudang(informasi_pickup_id);
CREATE INDEX idx_jadwal_gudang_hari ON jadwal_gudang(hari);

CREATE TRIGGER trg_jadwal_gudang_updated_at
    BEFORE UPDATE ON jadwal_gudang
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Seed default jadwal (Senin-Sabtu buka 09:00-18:00, Minggu tutup)
DO $$
DECLARE
    v_pickup_id UUID;
BEGIN
    SELECT id INTO v_pickup_id FROM informasi_pickup LIMIT 1;
    
    -- Minggu (0) - Tutup
    INSERT INTO jadwal_gudang (informasi_pickup_id, hari, jam_buka, jam_tutup, is_buka)
    VALUES (v_pickup_id, 0, NULL, NULL, false);
    
    -- Senin-Sabtu (1-6) - Buka 09:00-18:00
    INSERT INTO jadwal_gudang (informasi_pickup_id, hari, jam_buka, jam_tutup, is_buka)
    VALUES 
        (v_pickup_id, 1, '09:00', '18:00', true),
        (v_pickup_id, 2, '09:00', '18:00', true),
        (v_pickup_id, 3, '09:00', '18:00', true),
        (v_pickup_id, 4, '09:00', '18:00', true),
        (v_pickup_id, 5, '09:00', '18:00', true),
        (v_pickup_id, 6, '09:00', '18:00', true);
END $$;

COMMENT ON TABLE jadwal_gudang IS 'Jadwal operasional gudang per hari. Hari: 0=Minggu, 1=Senin, ..., 6=Sabtu';
