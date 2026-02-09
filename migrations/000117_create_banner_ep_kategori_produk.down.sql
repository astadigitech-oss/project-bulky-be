-- Rollback: Restore tujuan column and migrate data back

-- Step 1: Add tujuan column back
ALTER TABLE banner_event_promo 
ADD COLUMN tujuan VARCHAR(1000);

-- Step 2: Migrate data back to comma-separated
UPDATE banner_event_promo b
SET tujuan = (
    SELECT string_agg(bek.kategori_produk_id::TEXT, ',')
    FROM banner_ep_kategori_produk bek
    WHERE bek.banner_id = b.id
)
WHERE EXISTS (
    SELECT 1 FROM banner_ep_kategori_produk bek WHERE bek.banner_id = b.id
);

-- Step 3: Drop pivot table
DROP TABLE IF EXISTS banner_ep_kategori_produk;
