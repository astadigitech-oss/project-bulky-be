-- =====================================================
-- DISABLE Slug Triggers - Handle di Go Code
-- =====================================================
-- Karena GORM tidak pass slug ke trigger, kita handle di Go
-- Disable semua trigger rewrite_slug_on_delete
-- =====================================================

-- Drop triggers dari semua tabel
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

-- Function tetap ada untuk backward compatibility, tapi tidak dipakai
COMMENT ON FUNCTION rewrite_slug_on_delete() IS 
'DEPRECATED: Slug rewrite now handled in Go code. Function kept for backward compatibility.';
