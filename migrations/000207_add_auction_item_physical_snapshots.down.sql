ALTER TABLE auction_batch_items
    DROP CONSTRAINT IF EXISTS chk_auction_batch_items_physical_snapshot,
    DROP COLUMN IF EXISTS unit_panjang_cm,
    DROP COLUMN IF EXISTS unit_lebar_cm,
    DROP COLUMN IF EXISTS unit_tinggi_cm,
    DROP COLUMN IF EXISTS unit_volume_m3,
    DROP COLUMN IF EXISTS unit_berat_kg;

ALTER TABLE auction_batches
    DROP CONSTRAINT IF EXISTS chk_auction_batches_physical_source,
    DROP COLUMN IF EXISTS physical_source,
    ALTER COLUMN volume_m3 TYPE NUMERIC(12,3);
