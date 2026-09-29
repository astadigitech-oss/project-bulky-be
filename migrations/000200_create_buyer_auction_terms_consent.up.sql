CREATE TABLE buyer_auction_terms_consent (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    buyer_id UUID NOT NULL REFERENCES buyer(id) ON DELETE RESTRICT,
    bid_id UUID NOT NULL REFERENCES auction_bids(id) ON DELETE RESTRICT,
    terms_document_id UUID NOT NULL REFERENCES dokumen_kebijakan(id) ON DELETE RESTRICT,
    locale VARCHAR(5) NOT NULL CHECK (locale IN ('id', 'en')),
    content_hash CHAR(64) NOT NULL,
    disetujui_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_buyer_auction_terms_consent_bid UNIQUE (bid_id)
);

CREATE INDEX idx_batc_buyer_id ON buyer_auction_terms_consent(buyer_id);
CREATE INDEX idx_batc_terms_document_id ON buyer_auction_terms_consent(terms_document_id);
CREATE INDEX idx_batc_disetujui_at ON buyer_auction_terms_consent(disetujui_at DESC);

COMMENT ON TABLE buyer_auction_terms_consent IS 'Audit trail persetujuan syarat dan ketentuan lelang buyer pada saat mengirim bid.';
COMMENT ON COLUMN buyer_auction_terms_consent.content_hash IS 'SHA-256 konten syarat lelang sesuai locale yang disetujui buyer.';
