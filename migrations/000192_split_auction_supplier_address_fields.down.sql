ALTER TABLE auction_batches
    DROP CONSTRAINT IF EXISTS chk_auction_batch_supplier_coordinates,
    DROP CONSTRAINT IF EXISTS chk_auction_batch_origin;

ALTER TABLE auction_batches
    DROP COLUMN IF EXISTS supplier_longitude,
    DROP COLUMN IF EXISTS supplier_latitude,
    DROP COLUMN IF EXISTS supplier_kode_pos,
    DROP COLUMN IF EXISTS supplier_kelurahan,
    DROP COLUMN IF EXISTS supplier_kecamatan,
    DROP COLUMN IF EXISTS supplier_provinsi;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'auction_batches' AND column_name = 'supplier_kota'
    ) THEN
        ALTER TABLE auction_batches RENAME COLUMN supplier_kota TO supplier_city;
    END IF;
END $$;

ALTER TABLE auction_batches
    ADD CONSTRAINT chk_auction_batch_origin CHECK (
        status = 'DRAFT'
        OR (origin_type = 'BULKY_WAREHOUSE' AND warehouse_id IS NOT NULL AND supplier_name IS NULL AND supplier_address IS NULL AND supplier_city IS NULL)
        OR (origin_type = 'SUPPLIER' AND warehouse_id IS NULL AND supplier_name IS NOT NULL AND supplier_address IS NOT NULL AND supplier_city IS NOT NULL)
    );
