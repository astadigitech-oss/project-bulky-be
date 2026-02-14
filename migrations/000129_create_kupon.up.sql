CREATE TABLE kupon (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kode VARCHAR(50) NOT NULL,
    nama VARCHAR(255),
    deskripsi TEXT,
    jenis_diskon VARCHAR(20) NOT NULL CHECK (jenis_diskon IN ('persentase', 'jumlah_tetap')),
    nilai_diskon DECIMAL(15,2) NOT NULL CHECK (nilai_diskon > 0),
    minimal_pembelian DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (minimal_pembelian >= 0),
    limit_pemakaian INTEGER CHECK (limit_pemakaian IS NULL OR limit_pemakaian > 0),
    tanggal_kedaluarsa DATE NOT NULL,
    is_all_kategori BOOLEAN NOT NULL DEFAULT true,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT kupon_nilai_persentase_check CHECK (
        jenis_diskon != 'persentase' OR (nilai_diskon > 0 AND nilai_diskon <= 100)
    )
);

-- Unique constraint case-insensitive
CREATE UNIQUE INDEX idx_kupon_kode_unique ON kupon (LOWER(kode)) WHERE deleted_at IS NULL;

-- Standard indexes
CREATE INDEX idx_kupon_is_active ON kupon(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_kupon_tanggal_kedaluarsa ON kupon(tanggal_kedaluarsa) WHERE deleted_at IS NULL;
CREATE INDEX idx_kupon_created_at ON kupon(created_at DESC);

-- Trigger for updated_at
CREATE TRIGGER trg_kupon_updated_at
    BEFORE UPDATE ON kupon
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE kupon IS 'Tabel kupon/promo code untuk checkout';
COMMENT ON COLUMN kupon.kode IS 'Kode unik kupon (case-insensitive)';
COMMENT ON COLUMN kupon.jenis_diskon IS 'Jenis diskon: persentase atau jumlah_tetap';
COMMENT ON COLUMN kupon.nilai_diskon IS 'Nilai diskon sesuai jenis (max 100 jika persentase)';
COMMENT ON COLUMN kupon.minimal_pembelian IS 'Minimal harga per produk untuk bisa pakai kupon';
COMMENT ON COLUMN kupon.limit_pemakaian IS 'Batas total penggunaan kupon (NULL = unlimited)';
COMMENT ON COLUMN kupon.tanggal_kedaluarsa IS 'Tanggal berakhirnya kupon';
COMMENT ON COLUMN kupon.is_all_kategori IS 'true = berlaku semua kategori, false = cek tabel pivot';
