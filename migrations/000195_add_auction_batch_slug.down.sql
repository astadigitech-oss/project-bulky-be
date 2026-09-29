DROP INDEX IF EXISTS uq_auction_batches_slug_id;
DROP INDEX IF EXISTS uq_auction_batches_slug_en;

ALTER TABLE auction_batches
    DROP COLUMN IF EXISTS slug_id,
    DROP COLUMN IF EXISTS slug_en;
