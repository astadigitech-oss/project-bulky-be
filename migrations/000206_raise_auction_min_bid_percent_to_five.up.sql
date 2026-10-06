ALTER TABLE auction_batches
    ALTER COLUMN min_bid_percent SET DEFAULT 5;

UPDATE auction_batches
SET min_bid_percent = 5
WHERE status IN ('DRAFT', 'OPEN');
