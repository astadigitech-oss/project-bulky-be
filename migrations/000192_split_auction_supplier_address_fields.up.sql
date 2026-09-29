-- Memecah alamat gudang supplier menjadi komponen wilayah terstruktur dan koordinat
-- untuk kebutuhan kalkulasi estimasi ongkir Deliveree & Forwarder di Storefront.

ALTER TABLE auction_batches
    DROP CONSTRAINT IF EXISTS chk_auction_batch_origin;

-- Rename supplier_city ke supplier_kota jika supplier_city ada
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'auction_batches' AND column_name = 'supplier_city'
    ) THEN
        ALTER TABLE auction_batches RENAME COLUMN supplier_city TO supplier_kota;
    ELSE
        ALTER TABLE auction_batches ADD COLUMN IF NOT EXISTS supplier_kota VARCHAR(100);
    END IF;
END $$;

ALTER TABLE auction_batches
    ADD COLUMN IF NOT EXISTS supplier_provinsi VARCHAR(100),
    ADD COLUMN IF NOT EXISTS supplier_kecamatan VARCHAR(100),
    ADD COLUMN IF NOT EXISTS supplier_kelurahan VARCHAR(100),
    ADD COLUMN IF NOT EXISTS supplier_kode_pos VARCHAR(10),
    ADD COLUMN IF NOT EXISTS supplier_latitude NUMERIC(10,8),
    ADD COLUMN IF NOT EXISTS supplier_longitude NUMERIC(11,8);

ALTER TABLE auction_batches
    ADD CONSTRAINT chk_auction_batch_origin CHECK (
        status = 'DRAFT'
        OR (
            origin_type = 'BULKY_WAREHOUSE'
            AND warehouse_id IS NOT NULL
            AND supplier_name IS NULL
            AND supplier_address IS NULL
            AND supplier_provinsi IS NULL
            AND supplier_kota IS NULL
            AND supplier_kecamatan IS NULL
            AND supplier_kelurahan IS NULL
            AND supplier_kode_pos IS NULL
            AND supplier_latitude IS NULL
            AND supplier_longitude IS NULL
        )
        OR (
            origin_type = 'SUPPLIER'
            AND warehouse_id IS NULL
            AND supplier_name IS NOT NULL
            AND supplier_address IS NOT NULL
            AND supplier_provinsi IS NOT NULL
            AND supplier_kota IS NOT NULL
            AND supplier_kecamatan IS NOT NULL
            AND supplier_latitude IS NOT NULL
            AND supplier_longitude IS NOT NULL
        )
    ),
    ADD CONSTRAINT chk_auction_batch_supplier_coordinates CHECK (
        (supplier_latitude IS NULL AND supplier_longitude IS NULL)
        OR (
            supplier_latitude BETWEEN -90 AND 90
            AND supplier_longitude BETWEEN -180 AND 180
        )
    );
