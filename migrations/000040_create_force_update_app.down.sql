-- migrations/000040_create_force_update_app.down.sql

DROP TRIGGER IF EXISTS trg_single_active_force_update ON force_update_app;
DROP TRIGGER IF EXISTS trg_force_update_updated_at ON force_update_app;
DROP FUNCTION IF EXISTS fn_ensure_single_active_force_update();
DROP TABLE IF EXISTS force_update_app;
DROP TYPE IF EXISTS update_type;
