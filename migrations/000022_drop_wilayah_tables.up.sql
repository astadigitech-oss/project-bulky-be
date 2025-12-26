-- =====================================================
-- DROP WILAYAH TABLES
-- =====================================================
-- Rollback tabel master wilayah karena sudah tidak digunakan
-- Data wilayah akan di-fetch dari Google Maps API
-- 
-- Lihat: 12-REVISI-ALAMAT-BUYER.md
-- =====================================================

-- Drop alamat_buyer dulu karena ada FK ke kelurahan
DROP TABLE IF EXISTS alamat_buyer CASCADE;

-- Drop tables (urutan: child dulu, parent terakhir)
DROP TABLE IF EXISTS kelurahan CASCADE;
DROP TABLE IF EXISTS kecamatan CASCADE;
DROP TABLE IF EXISTS kota CASCADE;
DROP TABLE IF EXISTS provinsi CASCADE;
