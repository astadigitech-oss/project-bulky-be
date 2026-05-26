ALTER TABLE produk
    ADD COLUMN IF NOT EXISTS discrepancy_percentage DECIMAL(5,2) NOT NULL DEFAULT 0
        CHECK (discrepancy_percentage >= 0 AND discrepancy_percentage <= 100);

COMMENT ON COLUMN produk.discrepancy IS 'Deskripsi discrepancy/kekurangan produk';
COMMENT ON COLUMN produk.discrepancy_percentage IS 'Persentase discrepancy produk (0.00 - 100.00)';
