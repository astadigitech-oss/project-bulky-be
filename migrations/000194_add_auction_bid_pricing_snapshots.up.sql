ALTER TABLE auction_shipping_quotes
    ALTER COLUMN buyer_id DROP NOT NULL;

ALTER TABLE auction_bids
    ADD COLUMN shipping_provider_snapshot VARCHAR(50) NOT NULL DEFAULT '',
    ADD COLUMN shipping_service_snapshot VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN shipping_amount_snapshot NUMERIC(18,0) NOT NULL DEFAULT 0 CHECK (shipping_amount_snapshot >= 0),
    ADD COLUMN ppn_rate_snapshot NUMERIC(5,2) NOT NULL DEFAULT 0 CHECK (ppn_rate_snapshot >= 0 AND ppn_rate_snapshot <= 100),
    ADD COLUMN ppn_amount_snapshot NUMERIC(18,0) NOT NULL DEFAULT 0 CHECK (ppn_amount_snapshot >= 0),
    ADD COLUMN estimated_total_snapshot NUMERIC(18,0) NOT NULL DEFAULT 0 CHECK (estimated_total_snapshot >= 0);

CREATE INDEX idx_auction_shipping_quotes_batch_expiry
    ON auction_shipping_quotes (batch_id, expires_at);
