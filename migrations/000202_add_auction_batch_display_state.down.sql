DROP INDEX IF EXISTS idx_auction_batches_displayed_status_opened_at;

ALTER TABLE auction_batches
    DROP COLUMN IF EXISTS is_displayed;
