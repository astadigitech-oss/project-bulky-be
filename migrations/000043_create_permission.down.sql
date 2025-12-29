-- migrations/000043_create_permission.down.sql

DROP TRIGGER IF EXISTS trg_permission_updated_at ON permission;
DROP TABLE IF EXISTS permission;
