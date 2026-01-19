-- =====================================================
-- ROLLBACK: Remove triggers and function
-- =====================================================

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_kategori_produk ON kategori_produk;
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_merek_produk ON merek_produk;
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_kondisi_produk ON kondisi_produk;
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_kondisi_paket ON kondisi_paket;
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_sumber_produk ON sumber_produk;
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_warehouse ON warehouse;
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_tipe_produk ON tipe_produk;
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_produk ON produk;
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_dokumen_kebijakan ON dokumen_kebijakan;
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_disclaimer ON disclaimer;

-- Drop function
DROP FUNCTION IF EXISTS rewrite_slug_on_delete();
