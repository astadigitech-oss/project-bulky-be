-- migrations/000166_alter_mode_maintenance_dual_bahasa.up.sql

-- =====================================================
-- ALTER TABLE: mode_maintenance
-- =====================================================
-- Tambah dukungan dual bahasa (judul_en, deskripsi_en)
-- Scope tetap global (tanpa platform), single-active trigger tidak berubah
-- =====================================================

ALTER TABLE mode_maintenance
ADD COLUMN judul_en VARCHAR(100) NOT NULL DEFAULT '',
ADD COLUMN deskripsi_en TEXT NOT NULL DEFAULT '';

ALTER TABLE mode_maintenance
ALTER COLUMN judul_en DROP DEFAULT,
ALTER COLUMN deskripsi_en DROP DEFAULT;

COMMENT ON COLUMN mode_maintenance.judul_en IS 'Judul maintenance (Bahasa Inggris)';
COMMENT ON COLUMN mode_maintenance.deskripsi_en IS 'Deskripsi maintenance (Bahasa Inggris)';
