-- Auth Performance Indexes
-- Optimasi untuk endpoint authentication dan query lookup

-- Admin table indexes
CREATE INDEX IF NOT EXISTS idx_admin_email_active 
    ON admin(email) 
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_admin_role_id 
    ON admin(role_id);

CREATE INDEX IF NOT EXISTS idx_admin_is_active_deleted 
    ON admin(is_active) 
    WHERE deleted_at IS NULL;

-- Buyer table indexes
CREATE INDEX IF NOT EXISTS idx_buyer_email_active 
    ON buyer(email) 
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_buyer_username_active 
    ON buyer(username) 
    WHERE deleted_at IS NULL;

-- Admin Session table indexes
CREATE INDEX IF NOT EXISTS idx_admin_session_admin_id 
    ON admin_session(admin_id);

CREATE INDEX IF NOT EXISTS idx_admin_session_expires 
    ON admin_session(expires_at) 
    WHERE revoked_at IS NULL;

-- Role Permission lookup
CREATE INDEX IF NOT EXISTS idx_role_permission_role 
    ON role_permission(role_id);

CREATE INDEX IF NOT EXISTS idx_role_permission_permission 
    ON role_permission(permission_id);

-- Role table
CREATE INDEX IF NOT EXISTS idx_role_is_active 
    ON role(is_active);
