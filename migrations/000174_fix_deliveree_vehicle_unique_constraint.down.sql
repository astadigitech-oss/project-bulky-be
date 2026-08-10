-- migrations/000174_fix_deliveree_vehicle_unique_constraint.down.sql
-- Rollback: kembalikan ke partial unique index (kondisi awal migration 000172).

ALTER TABLE deliveree_vehicle_type
    DROP CONSTRAINT IF EXISTS uq_deliveree_vehicle_type_id_env;

CREATE UNIQUE INDEX idx_deliveree_vehicle_type_id_env
    ON deliveree_vehicle_type(id_deliveree, environment)
    WHERE deleted_at IS NULL;
