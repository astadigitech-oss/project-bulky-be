ALTER TABLE auction_batches
    ADD COLUMN physical_source VARCHAR(20) NOT NULL DEFAULT 'MANUAL',
    ADD CONSTRAINT chk_auction_batches_physical_source
        CHECK (physical_source IN ('MANUAL', 'ITEM_AGGREGATE')),
    ALTER COLUMN volume_m3 TYPE NUMERIC(12,6);

ALTER TABLE auction_batch_items
    ADD COLUMN unit_panjang_cm NUMERIC(12,3),
    ADD COLUMN unit_lebar_cm NUMERIC(12,3),
    ADD COLUMN unit_tinggi_cm NUMERIC(12,3),
    ADD COLUMN unit_volume_m3 NUMERIC(15,6),
    ADD COLUMN unit_berat_kg NUMERIC(12,3),
    ADD CONSTRAINT chk_auction_batch_items_physical_snapshot CHECK (
        (unit_panjang_cm IS NULL AND unit_lebar_cm IS NULL AND unit_tinggi_cm IS NULL AND unit_volume_m3 IS NULL AND unit_berat_kg IS NULL)
        OR
        (unit_panjang_cm IS NOT NULL AND unit_panjang_cm > 0
            AND unit_lebar_cm IS NOT NULL AND unit_lebar_cm > 0
            AND unit_tinggi_cm IS NOT NULL AND unit_tinggi_cm > 0
            AND unit_volume_m3 IS NOT NULL AND unit_volume_m3 > 0
            AND unit_berat_kg IS NOT NULL AND unit_berat_kg > 0)
    );
