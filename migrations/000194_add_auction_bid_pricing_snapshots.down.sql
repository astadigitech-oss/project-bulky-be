DROP INDEX IF EXISTS idx_auction_shipping_quotes_batch_expiry;

ALTER TABLE auction_bids
    DROP COLUMN IF EXISTS estimated_total_snapshot,
    DROP COLUMN IF EXISTS ppn_amount_snapshot,
    DROP COLUMN IF EXISTS ppn_rate_snapshot,
    DROP COLUMN IF EXISTS shipping_amount_snapshot,
    DROP COLUMN IF EXISTS shipping_service_snapshot,
    DROP COLUMN IF EXISTS shipping_provider_snapshot;

-- buyer_id intentionally remains nullable: quote rows may have been created
-- anonymously by the public estimate endpoint and must not be discarded.
