-- Rollback: buat ulang table refresh_token
CREATE TABLE IF NOT EXISTS refresh_token (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_type VARCHAR(10) NOT NULL CHECK (user_type IN ('ADMIN', 'BUYER')),
    user_id UUID NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    device_info TEXT,
    ip_address VARCHAR(45),
    expired_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_refresh_token_user ON refresh_token(user_type, user_id);
CREATE INDEX idx_refresh_token_token ON refresh_token(token);
CREATE INDEX idx_refresh_token_expired ON refresh_token(expired_at);
CREATE INDEX idx_refresh_token_is_revoked ON refresh_token(is_revoked);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_refresh_token_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_refresh_token_updated_at
    BEFORE UPDATE ON refresh_token
    FOR EACH ROW
    EXECUTE FUNCTION update_refresh_token_updated_at();
