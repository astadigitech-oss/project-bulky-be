-- Add dimensi & berat fields to produk table
ALTER TABLE produk
ADD COLUMN panjang DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (panjang >= 0),
ADD COLUMN lebar DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (lebar >= 0),
ADD COLUMN tinggi DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (tinggi >= 0),
ADD COLUMN berat DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (berat >= 0);

-- Add comments for documentation
COMMENT ON COLUMN produk.panjang IS 'Panjang produk dalam centimeter (cm)';
COMMENT ON COLUMN produk.lebar IS 'Lebar produk dalam centimeter (cm)';
COMMENT ON COLUMN produk.tinggi IS 'Tinggi produk dalam centimeter (cm)';
COMMENT ON COLUMN produk.berat IS 'Berat produk dalam kilogram (kg)';

-- Optional: Index untuk query shipping (jika sering filter by berat)
CREATE INDEX idx_produk_berat ON produk(berat) WHERE deleted_at IS NULL;
