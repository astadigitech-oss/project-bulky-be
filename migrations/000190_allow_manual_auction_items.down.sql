ALTER TABLE auction_batch_items
    DROP CONSTRAINT IF EXISTS chk_auction_batch_item_source,
    DROP COLUMN IF EXISTS source_type,
    ALTER COLUMN produk_id SET NOT NULL;

ALTER TABLE auction_batches
    DROP CONSTRAINT IF EXISTS chk_auction_batch_origin,
    DROP COLUMN IF EXISTS supplier_city,
    DROP COLUMN IF EXISTS supplier_address,
    DROP COLUMN IF EXISTS supplier_name,
    DROP COLUMN IF EXISTS origin_type;
