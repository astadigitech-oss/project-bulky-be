-- Drop performance indexes

DROP INDEX IF EXISTS idx_admin_email_active;
DROP INDEX IF EXISTS idx_admin_role_id;
DROP INDEX IF EXISTS idx_admin_is_active_deleted;

DROP INDEX IF EXISTS idx_buyer_email_active;
DROP INDEX IF EXISTS idx_buyer_username_active;

DROP INDEX IF EXISTS idx_admin_session_admin_id;
DROP INDEX IF EXISTS idx_admin_session_expires;

DROP INDEX IF EXISTS idx_role_permission_role;
DROP INDEX IF EXISTS idx_role_permission_permission;

DROP INDEX IF EXISTS idx_role_is_active;
