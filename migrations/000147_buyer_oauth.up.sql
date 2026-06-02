-- migrations/000147_buyer_oauth.up.sql

-- =====================================================
-- STEP 1: Buat tabel buyer_oauth
-- =====================================================

CREATE TABLE buyer_oauth (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_id     UUID        NOT NULL REFERENCES buyer(id) ON DELETE CASCADE,
    provider     VARCHAR(50) NOT NULL,
    provider_uid VARCHAR(255) NOT NULL,
    email        VARCHAR(255),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

-- =====================================================
-- STEP 2: Index
-- =====================================================

-- Index untuk lookup buyer_id (semua OAuth method milik 1 buyer)
CREATE INDEX idx_buyer_oauth_buyer_id
    ON buyer_oauth(buyer_id)
    WHERE deleted_at IS NULL;

-- Unique constraint: satu provider_uid hanya bisa linked ke 1 buyer
CREATE UNIQUE INDEX buyer_oauth_provider_uid_unique
    ON buyer_oauth(provider, provider_uid)
    WHERE deleted_at IS NULL;

-- =====================================================
-- STEP 3: Comments
-- =====================================================

COMMENT ON TABLE buyer_oauth IS 'Relasi buyer dengan OAuth provider (Google, Apple). Satu buyer bisa punya lebih dari satu OAuth provider.';
COMMENT ON COLUMN buyer_oauth.provider IS 'Nama OAuth provider: ''google'' atau ''apple''';
COMMENT ON COLUMN buyer_oauth.provider_uid IS 'ID unik dari provider (Google: sub claim JWT, Apple: user field)';
COMMENT ON COLUMN buyer_oauth.email IS 'Email dari provider saat OAuth. Apple mungkin mengembalikan relay email (xxx@privaterelay.appleid.com).';
