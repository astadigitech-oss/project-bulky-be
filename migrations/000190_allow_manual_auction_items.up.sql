-- Batch lelang tidak harus berasal dari katalog Bulky. Produk supplier dapat
-- dicatat sebagai snapshot MANUAL tanpa FK produk dan tanpa reservasi stok.
ALTER TABLE auction_batches
    ADD COLUMN origin_type VARCHAR(20) NOT NULL DEFAULT 'BULKY_WAREHOUSE'
        CHECK (origin_type IN ('BULKY_WAREHOUSE', 'SUPPLIER')),
    ADD COLUMN supplier_name VARCHAR(255),
    ADD COLUMN supplier_address TEXT,
    ADD COLUMN supplier_city VARCHAR(100),
    ADD CONSTRAINT chk_auction_batch_origin CHECK (
        status = 'DRAFT'
        OR (origin_type = 'BULKY_WAREHOUSE' AND warehouse_id IS NOT NULL AND supplier_name IS NULL AND supplier_address IS NULL AND supplier_city IS NULL)
        OR
        (origin_type = 'SUPPLIER' AND warehouse_id IS NULL AND supplier_name IS NOT NULL AND supplier_address IS NOT NULL AND supplier_city IS NOT NULL)
    );

ALTER TABLE auction_batch_items
    ALTER COLUMN produk_id DROP NOT NULL,
    ADD COLUMN source_type VARCHAR(10) NOT NULL DEFAULT 'CATALOG'
        CHECK (source_type IN ('CATALOG', 'MANUAL')),
    ADD CONSTRAINT chk_auction_batch_item_source CHECK (
        (source_type = 'CATALOG' AND produk_id IS NOT NULL)
        OR
        (source_type = 'MANUAL' AND produk_id IS NULL)
    );
