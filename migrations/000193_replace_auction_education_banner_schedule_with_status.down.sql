DROP INDEX IF EXISTS idx_auction_education_banners_published;

ALTER TABLE auction_education_banners
    DROP COLUMN is_published,
    ADD COLUMN tanggal_mulai TIMESTAMPTZ NULL,
    ADD COLUMN tanggal_selesai TIMESTAMPTZ NULL,
    ADD CONSTRAINT chk_auction_education_banner_schedule CHECK (
        tanggal_mulai IS NULL OR tanggal_selesai IS NULL OR tanggal_selesai > tanggal_mulai
    );

CREATE INDEX idx_auction_education_banners_schedule
    ON auction_education_banners (tanggal_mulai, tanggal_selesai)
    WHERE deleted_at IS NULL;
