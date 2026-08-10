-- migrations/000172_create_deliveree_vehicle_type.up.sql
-- Master data kendaraan Deliveree, ditarik via tombol Sync dari
-- GET /public_api/v10/vehicle_types (sandbox & production terpisah).
-- Dipakai sebagai basis keputusan pemilihan kendaraan saat create booking
-- Deliveree (berdasarkan kubikasi & berat total pesanan, menggantikan logic
-- lama yang berdasarkan jumlah palet/qty).

CREATE TABLE deliveree_vehicle_type (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    id_deliveree INTEGER NOT NULL,
    environment VARCHAR(20) NOT NULL CHECK (environment IN ('sandbox', 'production')),
    kubikasi_max DECIMAL(10,3) NOT NULL DEFAULT 0,
    berat_max DECIMAL(10,2) NOT NULL DEFAULT 0,
    threshold_kubikasi DECIMAL(10,3) NOT NULL DEFAULT 0,
    threshold_berat DECIMAL(10,2) NOT NULL DEFAULT 0,
    cargo_length DECIMAL(10,2),
    cargo_width DECIMAL(10,2),
    cargo_height DECIMAL(10,2),
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Satu id_deliveree hanya boleh muncul sekali per environment (dipakai untuk upsert saat Sync)
CREATE UNIQUE INDEX idx_deliveree_vehicle_type_id_env
    ON deliveree_vehicle_type(id_deliveree, environment)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_deliveree_vehicle_type_env_active
    ON deliveree_vehicle_type(environment, is_active)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_deliveree_vehicle_type_kubikasi
    ON deliveree_vehicle_type(environment, kubikasi_max)
    WHERE deleted_at IS NULL AND is_active = true;

CREATE TRIGGER trg_deliveree_vehicle_type_updated_at
    BEFORE UPDATE ON deliveree_vehicle_type
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE deliveree_vehicle_type IS 'Master data kendaraan Deliveree (sandbox & production), ditarik via Sync dari API vehicle_types. Dipakai untuk logic pemilihan kendaraan saat create booking Deliveree berdasarkan kubikasi & berat.';
COMMENT ON COLUMN deliveree_vehicle_type.id_deliveree IS 'vehicle_type_id di sisi Deliveree, dipakai sebagai vehicle_type_id saat create booking';
COMMENT ON COLUMN deliveree_vehicle_type.environment IS 'sandbox | production - vehicle_type_id berbeda antar environment';
COMMENT ON COLUMN deliveree_vehicle_type.kubikasi_max IS 'Kapasitas kubikasi maksimal (m3), dari field cargo_cubic_meter Deliveree';
COMMENT ON COLUMN deliveree_vehicle_type.berat_max IS 'Kapasitas berat maksimal (kg), dari field cargo_weight Deliveree';
COMMENT ON COLUMN deliveree_vehicle_type.threshold_kubikasi IS 'Auto-computed saat Sync = kubikasi_max kendaraan satu tingkat di bawahnya (0 untuk kendaraan terkecil). Bisa di-override manual oleh Super Admin.';
COMMENT ON COLUMN deliveree_vehicle_type.threshold_berat IS 'Auto-computed saat Sync = berat_max kendaraan satu tingkat di bawahnya (0 untuk kendaraan terkecil). Bisa di-override manual oleh Super Admin.';
COMMENT ON COLUMN deliveree_vehicle_type.is_active IS 'false jika kendaraan sudah tidak muncul lagi di response API Deliveree saat Sync terakhir (soft-disable), atau dinonaktifkan manual';
COMMENT ON COLUMN deliveree_vehicle_type.last_synced_at IS 'Waktu terakhir data ini berhasil disinkronkan dari API Deliveree';
