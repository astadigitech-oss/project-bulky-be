CREATE TABLE admin_session (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id UUID NOT NULL REFERENCES admin(id) ON DELETE CASCADE,
    token VARCHAR(500) NOT NULL UNIQUE,
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_admin_session_admin_id ON admin_session(admin_id);
CREATE INDEX idx_admin_session_token ON admin_session(token);
CREATE INDEX idx_admin_session_expires_at ON admin_session(expires_at);

-- Function to clean expired sessions
CREATE OR REPLACE FUNCTION clean_expired_sessions()
RETURNS void AS $$
BEGIN
    DELETE FROM admin_session WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Table & Column Comments
COMMENT ON TABLE admin_session IS 'Menyimpan session/refresh token admin';
COMMENT ON COLUMN admin_session.id IS 'Primary key UUID';
COMMENT ON COLUMN admin_session.admin_id IS 'FK ke admin';
COMMENT ON COLUMN admin_session.token IS 'Refresh token (unique)';
COMMENT ON COLUMN admin_session.ip_address IS 'IP address saat login';
COMMENT ON COLUMN admin_session.user_agent IS 'Browser/device info';
COMMENT ON COLUMN admin_session.expires_at IS 'Waktu expired token';
COMMENT ON COLUMN admin_session.created_at IS 'Waktu dibuat';
