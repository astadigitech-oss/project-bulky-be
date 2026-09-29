ALTER TABLE auction_education_banners
    DROP CONSTRAINT IF EXISTS chk_auction_education_banner_schedule,
    DROP COLUMN tanggal_mulai,
    DROP COLUMN tanggal_selesai,
    ADD COLUMN is_published BOOLEAN NOT NULL DEFAULT FALSE;

DROP INDEX IF EXISTS idx_auction_education_banners_schedule;
CREATE INDEX idx_auction_education_banners_published
    ON auction_education_banners (is_published)
    WHERE deleted_at IS NULL;
