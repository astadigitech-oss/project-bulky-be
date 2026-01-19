-- =====================================================
-- TABEL: formulir_partai_besar_config
-- =====================================================
-- Konfigurasi formulir partai besar (singleton)
-- Menyimpan daftar email penerima notifikasi
-- =====================================================

CREATE TABLE formulir_partai_besar_config (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    daftar_email TEXT NOT NULL DEFAULT '[]',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trg_formulir_config_updated_at
    BEFORE UPDATE ON formulir_partai_besar_config
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Singleton constraint
CREATE OR REPLACE FUNCTION check_formulir_config_singleton()
RETURNS TRIGGER AS $$
DECLARE
    v_count INT;
BEGIN
    SELECT COUNT(*) INTO v_count FROM formulir_partai_besar_config;
    IF v_count >= 1 THEN
        RAISE EXCEPTION 'Hanya boleh ada 1 konfigurasi formulir (gunakan UPDATE)';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_formulir_config_singleton
    BEFORE INSERT ON formulir_partai_besar_config
    FOR EACH ROW
    EXECUTE FUNCTION check_formulir_config_singleton();

-- Seed default config
INSERT INTO formulir_partai_besar_config (daftar_email) VALUES ('[]');

COMMENT ON TABLE formulir_partai_besar_config IS 'Konfigurasi formulir partai besar. Singleton.';
COMMENT ON COLUMN formulir_partai_besar_config.daftar_email IS 'JSON array email penerima notifikasi';
