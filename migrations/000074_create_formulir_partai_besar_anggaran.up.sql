-- =====================================================
-- TABEL: formulir_partai_besar_anggaran
-- =====================================================
-- Opsi anggaran/budget untuk formulir partai besar
-- =====================================================

CREATE TABLE formulir_partai_besar_anggaran (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    label VARCHAR(100) NOT NULL,
    urutan INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_formulir_anggaran_urutan ON formulir_partai_besar_anggaran(urutan);
CREATE INDEX idx_formulir_anggaran_deleted_at ON formulir_partai_besar_anggaran(deleted_at);

CREATE TRIGGER trg_formulir_anggaran_updated_at
    BEFORE UPDATE ON formulir_partai_besar_anggaran
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Seed default anggaran
INSERT INTO formulir_partai_besar_anggaran (label, urutan) VALUES
    ('Rp 25.000.000 - Rp 50.000.000', 1),
    ('Rp 51.000.000 - Rp 100.000.000', 2),
    ('Rp 101.000.000 - Rp 200.000.000', 3),
    ('> Rp 200.000.000', 4);

COMMENT ON TABLE formulir_partai_besar_anggaran IS 'Opsi anggaran/budget untuk formulir partai besar';
