CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE seasonal_campaign (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT false,
    tanggal_mulai TIMESTAMPTZ,
    tanggal_selesai TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    web_logo_url VARCHAR(500),
    web_navbar_decoration_url VARCHAR(500),
    mobile_top_app_bar_ornament_url VARCHAR(500),
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_seasonal_campaign_period CHECK (
        (tanggal_mulai IS NULL AND tanggal_selesai IS NULL) OR
        (tanggal_mulai IS NOT NULL AND tanggal_selesai IS NOT NULL AND tanggal_selesai > tanggal_mulai)
    ),
    CONSTRAINT chk_seasonal_campaign_published_period CHECK (
        NOT is_published OR (tanggal_mulai IS NOT NULL AND tanggal_selesai IS NOT NULL)
    )
);

CREATE INDEX idx_seasonal_campaign_period ON seasonal_campaign (tanggal_mulai, tanggal_selesai) WHERE deleted_at IS NULL;
CREATE INDEX idx_seasonal_campaign_published ON seasonal_campaign (is_published) WHERE deleted_at IS NULL;

ALTER TABLE seasonal_campaign ADD CONSTRAINT seasonal_campaign_no_published_overlap
EXCLUDE USING gist (tstzrange(tanggal_mulai, tanggal_selesai, '[]') WITH &&)
WHERE (is_published AND cancelled_at IS NULL AND deleted_at IS NULL);

CREATE TRIGGER trg_seasonal_campaign_updated_at
BEFORE UPDATE ON seasonal_campaign
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE seasonal_campaign IS 'Konfigurasi branding campaign seasonal yang dikelola Admin Backend';
