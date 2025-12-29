-- migrations/000046_create_refresh_token.up.sql

CREATE TABLE refresh_token (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_type VARCHAR(20) NOT NULL CHECK (user_type IN ('ADMIN', 'BUYER')),
    user_id UUID NOT NULL,
    token VARCHAR(500) NOT NULL UNIQUE,
    device_info VARCHAR(255),
    ip_address VARCHAR(50),
    expired_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT false,
    revoked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_refresh_token_token ON refresh_token(token);
CREATE INDEX idx_refresh_token_user ON refresh_token(user_type, user_id);
CREATE INDEX idx_refresh_token_expired ON refresh_token(expired_at);
CREATE INDEX idx_refresh_token_active ON refresh_token(user_type, user_id, is_revoked, expired_at) 
    WHERE is_revoked = false;

COMMENT ON TABLE refresh_token IS 'Menyimpan refresh token untuk revocation & multi-device support';
COMMENT ON COLUMN refresh_token.token IS 'Hashed refresh token';
