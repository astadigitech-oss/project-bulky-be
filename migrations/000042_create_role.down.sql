-- migrations/000042_create_role.down.sql

DROP TRIGGER IF EXISTS trg_role_updated_at ON role;
DROP TABLE IF EXISTS role;
