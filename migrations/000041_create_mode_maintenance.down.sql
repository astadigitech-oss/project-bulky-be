-- migrations/000041_create_mode_maintenance.down.sql

DROP TRIGGER IF EXISTS trg_single_active_maintenance ON mode_maintenance;
DROP TRIGGER IF EXISTS trg_maintenance_updated_at ON mode_maintenance;
DROP FUNCTION IF EXISTS fn_ensure_single_active_maintenance();
DROP TABLE IF EXISTS mode_maintenance;
DROP TYPE IF EXISTS maintenance_type;
