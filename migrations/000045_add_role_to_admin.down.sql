-- migrations/000045_add_role_to_admin.down.sql

DROP INDEX IF EXISTS idx_admin_role_id;
ALTER TABLE admin DROP COLUMN IF EXISTS role_id;
