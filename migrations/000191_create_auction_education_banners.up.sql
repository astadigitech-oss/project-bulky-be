CREATE TABLE auction_education_banners (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    gambar_url_id VARCHAR(500) NOT NULL,
    gambar_url_en VARCHAR(500) NOT NULL,
    urutan INTEGER NOT NULL DEFAULT 0 CHECK (urutan >= 0),
    tanggal_mulai TIMESTAMPTZ NULL,
    tanggal_selesai TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT chk_auction_education_banner_schedule CHECK (
        tanggal_mulai IS NULL OR tanggal_selesai IS NULL OR tanggal_selesai > tanggal_mulai
    )
);

CREATE INDEX idx_auction_education_banners_order ON auction_education_banners (urutan ASC, id ASC) WHERE deleted_at IS NULL;
CREATE INDEX idx_auction_education_banners_schedule ON auction_education_banners (tanggal_mulai, tanggal_selesai) WHERE deleted_at IS NULL;

CREATE TRIGGER update_auction_education_banners_updated_at
    BEFORE UPDATE ON auction_education_banners
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
