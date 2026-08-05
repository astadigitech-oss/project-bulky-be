-- migrations/000163_create_buyer_disclaimer_consent.up.sql
-- Audit trail persetujuan disclaimer oleh buyer, dikorelasikan dengan pesanan

CREATE TABLE buyer_disclaimer_consent (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    buyer_id    UUID NOT NULL REFERENCES buyer(id),
    pesanan_id  UUID NOT NULL REFERENCES pesanan(id),
    disclaimer_id UUID NOT NULL REFERENCES disclaimer(id),
    disetujui_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address  VARCHAR(45),
    user_agent  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bdc_buyer_id    ON buyer_disclaimer_consent(buyer_id);
CREATE INDEX idx_bdc_pesanan_id  ON buyer_disclaimer_consent(pesanan_id);
CREATE INDEX idx_bdc_disclaimer_id ON buyer_disclaimer_consent(disclaimer_id);
CREATE INDEX idx_bdc_disetujui_at ON buyer_disclaimer_consent(disetujui_at DESC);

COMMENT ON TABLE buyer_disclaimer_consent IS 'Audit trail persetujuan disclaimer pembelian oleh buyer, dikorelasikan dengan pesanan.';
COMMENT ON COLUMN buyer_disclaimer_consent.ip_address IS 'IP address buyer saat menyetujui disclaimer (IPv4/IPv6).';
COMMENT ON COLUMN buyer_disclaimer_consent.user_agent IS 'User-Agent HTTP header dari client buyer.';
