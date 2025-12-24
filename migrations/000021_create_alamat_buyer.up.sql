CREATE TABLE alamat_buyer (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    buyer_id UUID NOT NULL REFERENCES buyer(id) ON DELETE CASCADE,
    kelurahan_id UUID NOT NULL REFERENCES kelurahan(id),
    label VARCHAR(50) NOT NULL,
    nama_penerima VARCHAR(100) NOT NULL,
    telepon_penerima VARCHAR(20) NOT NULL,
    kode_pos VARCHAR(10) NOT NULL,
    alamat_lengkap TEXT NOT NULL,
    catatan TEXT,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_alamat_buyer_buyer_id ON alamat_buyer(buyer_id);
CREATE INDEX idx_alamat_buyer_kelurahan_id ON alamat_buyer(kelurahan_id);
CREATE INDEX idx_alamat_buyer_is_default ON alamat_buyer(buyer_id, is_default) WHERE deleted_at IS NULL;

-- Trigger for updated_at
CREATE TRIGGER update_alamat_buyer_updated_at
    BEFORE UPDATE ON alamat_buyer
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to ensure only one default address per buyer
CREATE OR REPLACE FUNCTION ensure_single_default_alamat_buyer()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_default = true THEN
        UPDATE alamat_buyer 
        SET is_default = false 
        WHERE buyer_id = NEW.buyer_id 
          AND id != NEW.id 
          AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_single_default_alamat_buyer
    AFTER INSERT OR UPDATE ON alamat_buyer
    FOR EACH ROW
    WHEN (NEW.is_default = true)
    EXECUTE FUNCTION ensure_single_default_alamat_buyer();

-- Function to set first address as default
CREATE OR REPLACE FUNCTION set_first_alamat_buyer_as_default()
RETURNS TRIGGER AS $$
DECLARE
    address_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO address_count 
    FROM alamat_buyer 
    WHERE buyer_id = NEW.buyer_id AND deleted_at IS NULL;
    
    IF address_count = 0 THEN
        NEW.is_default := true;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_first_alamat_buyer_default
    BEFORE INSERT ON alamat_buyer
    FOR EACH ROW
    EXECUTE FUNCTION set_first_alamat_buyer_as_default();

-- Table & Column Comments
COMMENT ON TABLE alamat_buyer IS 'Menyimpan alamat pengiriman buyer. Satu buyer dapat memiliki banyak alamat, dengan satu alamat sebagai default.';
COMMENT ON COLUMN alamat_buyer.id IS 'Primary key UUID';
COMMENT ON COLUMN alamat_buyer.buyer_id IS 'Foreign key ke tabel buyer';
COMMENT ON COLUMN alamat_buyer.kelurahan_id IS 'Foreign key ke tabel kelurahan. Dari kelurahan bisa didapat kecamatan, kota, dan provinsi.';
COMMENT ON COLUMN alamat_buyer.label IS 'Label alamat seperti: Rumah, Kantor, Toko, dll';
COMMENT ON COLUMN alamat_buyer.nama_penerima IS 'Nama orang yang akan menerima paket';
COMMENT ON COLUMN alamat_buyer.telepon_penerima IS 'Nomor telepon penerima untuk dihubungi kurir';
COMMENT ON COLUMN alamat_buyer.kode_pos IS 'Kode pos alamat pengiriman';
COMMENT ON COLUMN alamat_buyer.alamat_lengkap IS 'Alamat lengkap termasuk nama jalan, nomor rumah, RT/RW';
COMMENT ON COLUMN alamat_buyer.catatan IS 'Catatan tambahan untuk kurir seperti patokan lokasi';
COMMENT ON COLUMN alamat_buyer.is_default IS 'Alamat default yang akan otomatis dipilih saat checkout. Hanya boleh ada 1 per buyer.';
COMMENT ON COLUMN alamat_buyer.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN alamat_buyer.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN alamat_buyer.deleted_at IS 'Soft delete timestamp';
