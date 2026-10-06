UPDATE auction_batches
SET min_bid_percent = 0.1
WHERE status IN ('DRAFT', 'OPEN');

ALTER TABLE auction_batches
    ALTER COLUMN min_bid_percent SET DEFAULT 0.1;
