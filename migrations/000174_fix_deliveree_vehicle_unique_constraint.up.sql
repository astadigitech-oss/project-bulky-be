-- migrations/000174_fix_deliveree_vehicle_unique_constraint.up.sql
-- Migration 000172 membuat partial unique index (WHERE deleted_at IS NULL) pada
-- (id_deliveree, environment). Partial index TIDAK bisa dipakai PostgreSQL untuk
-- klausa ON CONFLICT saat Upsert di Sync (error 42P10: no unique or exclusion
-- constraint matching the ON CONFLICT specification).
-- Fix: ganti dengan UNIQUE CONSTRAINT biasa pada (id_deliveree, environment).
-- Aman karena deaktivasi kendaraan memakai is_active=false (bukan soft-delete),
-- dan tabel ini tidak punya operasi hard-delete.

DROP INDEX IF EXISTS idx_deliveree_vehicle_type_id_env;

ALTER TABLE deliveree_vehicle_type
    ADD CONSTRAINT uq_deliveree_vehicle_type_id_env UNIQUE (id_deliveree, environment);
