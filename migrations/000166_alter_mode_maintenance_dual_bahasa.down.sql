-- migrations/000166_alter_mode_maintenance_dual_bahasa.down.sql

ALTER TABLE mode_maintenance
DROP COLUMN IF EXISTS judul_en,
DROP COLUMN IF EXISTS deskripsi_en;
