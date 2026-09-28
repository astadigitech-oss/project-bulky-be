ALTER TABLE auction_batches
    ADD COLUMN IF NOT EXISTS is_displayed BOOLEAN NOT NULL DEFAULT TRUE;

CREATE INDEX IF NOT EXISTS idx_auction_batches_displayed_status_opened_at
    ON auction_batches (status, opened_at DESC, id DESC)
    WHERE is_displayed = TRUE;

COMMENT ON COLUMN auction_batches.is_displayed IS
    'Menentukan apakah batch OPEN atau SOLD dapat tampil di Storefront; tidak mengubah status lelang atau kemampuan bid.';
