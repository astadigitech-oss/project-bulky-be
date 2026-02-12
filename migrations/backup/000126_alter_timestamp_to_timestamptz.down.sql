-- =====================================================
-- ROLLBACK: TIMESTAMPTZ to TIMESTAMP
-- =====================================================
-- WARNING: This will LOSE timezone information!
-- Data will be converted to Asia/Jakarta local time
-- and stored as TIMESTAMP (without timezone)
-- =====================================================

SET timezone = 'Asia/Jakarta';

-- =====================================================
-- DROP ALL TRIGGERS FIRST
-- =====================================================
-- Must drop all triggers that reference timestamp columns before ALTER TYPE

-- Master Data
DROP TRIGGER IF EXISTS update_kategori_produk_updated_at ON kategori_produk;
DROP TRIGGER IF EXISTS update_merek_produk_updated_at ON merek_produk;
DROP TRIGGER IF EXISTS update_kondisi_produk_updated_at ON kondisi_produk;
DROP TRIGGER IF EXISTS update_kondisi_paket_updated_at ON kondisi_paket;
DROP TRIGGER IF EXISTS update_sumber_produk_updated_at ON sumber_produk;

-- Produk Module
DROP TRIGGER IF EXISTS update_warehouse_updated_at ON warehouse;
DROP TRIGGER IF EXISTS update_tipe_produk_updated_at ON tipe_produk;
DROP TRIGGER IF EXISTS update_diskon_kategori_updated_at ON diskon_kategori;
DROP TRIGGER IF EXISTS update_banner_tipe_produk_updated_at ON banner_tipe_produk;
DROP TRIGGER IF EXISTS update_produk_updated_at ON produk;

-- Admin & Auth
DROP TRIGGER IF EXISTS update_admin_updated_at ON admin;
DROP TRIGGER IF EXISTS trg_role_updated_at ON role;
DROP TRIGGER IF EXISTS trg_permission_updated_at ON permission;

-- Buyer Module
DROP TRIGGER IF EXISTS update_buyer_updated_at ON buyer;
DROP TRIGGER IF EXISTS trg_alamat_buyer_updated_at ON alamat_buyer;
DROP TRIGGER IF EXISTS trg_single_default_alamat ON alamat_buyer;
DROP TRIGGER IF EXISTS trg_first_alamat_default ON alamat_buyer;
DROP TRIGGER IF EXISTS trg_prevent_delete_default ON alamat_buyer;

-- Wilayah (if exists)
DROP TRIGGER IF EXISTS update_provinsi_updated_at ON provinsi;
DROP TRIGGER IF EXISTS update_kota_updated_at ON kota;
DROP TRIGGER IF EXISTS update_kecamatan_updated_at ON kecamatan;
DROP TRIGGER IF EXISTS update_kelurahan_updated_at ON kelurahan;

-- Marketing
DROP TRIGGER IF EXISTS trg_hero_section_updated_at ON hero_section;
DROP TRIGGER IF EXISTS trg_banner_event_updated_at ON banner_event_promo;

-- Transaksi
DROP TRIGGER IF EXISTS trg_metode_pembayaran_group_updated_at ON metode_pembayaran_group;
DROP TRIGGER IF EXISTS trg_metode_pembayaran_updated_at ON metode_pembayaran;
DROP TRIGGER IF EXISTS trg_ppn_updated_at ON ppn;
DROP TRIGGER IF EXISTS trg_pesanan_updated_at ON pesanan;
DROP TRIGGER IF EXISTS trg_pesanan_item_updated_at ON pesanan_item;
DROP TRIGGER IF EXISTS trg_pesanan_pembayaran_updated_at ON pesanan_pembayaran;

-- Operasional
DROP TRIGGER IF EXISTS trg_informasi_pickup_updated_at ON informasi_pickup;
DROP TRIGGER IF EXISTS trg_jadwal_gudang_updated_at ON jadwal_gudang;
DROP TRIGGER IF EXISTS trg_dokumen_kebijakan_updated_at ON dokumen_kebijakan;
DROP TRIGGER IF EXISTS trg_disclaimer_updated_at ON disclaimer;

-- Ulasan
DROP TRIGGER IF EXISTS trg_ulasan_updated_at ON ulasan;

-- Formulir & Komunikasi
DROP TRIGGER IF EXISTS trg_formulir_config_updated_at ON formulir_partai_besar_config;
DROP TRIGGER IF EXISTS trg_formulir_anggaran_updated_at ON formulir_partai_besar_anggaran;
DROP TRIGGER IF EXISTS trg_whatsapp_handler_updated_at ON whatsapp_handler;

-- Sistem Kontrol
DROP TRIGGER IF EXISTS trg_force_update_updated_at ON force_update_app;
DROP TRIGGER IF EXISTS trg_maintenance_updated_at ON mode_maintenance;

-- Blog & Video
DROP TRIGGER IF EXISTS trg_kategori_blog_updated_at ON kategori_blog;
DROP TRIGGER IF EXISTS trg_label_blog_updated_at ON label_blog;
DROP TRIGGER IF EXISTS trg_blog_updated_at ON blog;
DROP TRIGGER IF EXISTS trg_kategori_video_updated_at ON kategori_video;
DROP TRIGGER IF EXISTS trg_video_updated_at ON video;

-- FAQ
DROP TRIGGER IF EXISTS trg_faq_updated_at ON faq;

-- Other business logic triggers that reference deleted_at or timestamps
DROP TRIGGER IF EXISTS trg_single_active_disclaimer ON disclaimer;
DROP TRIGGER IF EXISTS trg_single_active_wa_handler ON whatsapp_handler;
DROP TRIGGER IF EXISTS trg_single_active_maintenance ON mode_maintenance;
DROP TRIGGER IF EXISTS trg_single_active_force_update ON force_update_app;
DROP TRIGGER IF EXISTS trg_single_active_ppn ON ppn;
DROP TRIGGER IF EXISTS trg_single_default_hero ON hero_section;
DROP TRIGGER IF EXISTS trg_hero_section_auto_sync_insert ON hero_section;
DROP TRIGGER IF EXISTS trg_formulir_config_singleton ON formulir_partai_besar_config;
DROP TRIGGER IF EXISTS trg_informasi_pickup_singleton ON informasi_pickup;
DROP TRIGGER IF EXISTS trg_log_pesanan_status_change ON pesanan;
DROP TRIGGER IF EXISTS trg_calculate_subtotal ON pesanan_item;
DROP TRIGGER IF EXISTS trg_check_max_split_payment ON pesanan_pembayaran;
DROP TRIGGER IF EXISTS trg_validate_ulasan_order ON ulasan;
DROP TRIGGER IF EXISTS trg_set_ulasan_approved_at ON ulasan;
DROP TRIGGER IF EXISTS trg_blog_published_at ON blog;
DROP TRIGGER IF EXISTS trg_video_published_at ON video;
DROP TRIGGER IF EXISTS trg_generate_slug_disclaimer ON disclaimer;
DROP TRIGGER IF EXISTS trg_generate_slug_dokumen ON dokumen_kebijakan;
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

-- =====================================================
-- 1. MASTER DATA
-- =====================================================

-- kategori_produk
ALTER TABLE kategori_produk
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- merek_produk
ALTER TABLE merek_produk
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- kondisi_produk
ALTER TABLE kondisi_produk
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- kondisi_paket
ALTER TABLE kondisi_paket
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- sumber_produk
ALTER TABLE sumber_produk
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- 2. PRODUK MODULE
-- =====================================================

-- warehouse
ALTER TABLE warehouse
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- tipe_produk
ALTER TABLE tipe_produk
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- diskon_kategori
ALTER TABLE diskon_kategori
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- banner_tipe_produk
ALTER TABLE banner_tipe_produk
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- produk
ALTER TABLE produk
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- produk_gambar
ALTER TABLE produk_gambar
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta';

-- produk_dokumen
ALTER TABLE produk_dokumen
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- 3. ADMIN & AUTH
-- =====================================================

-- admin
ALTER TABLE admin
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN last_login_at TYPE TIMESTAMP USING last_login_at AT TIME ZONE 'Asia/Jakarta';

-- admin_session
ALTER TABLE admin_session
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN expires_at TYPE TIMESTAMP USING expires_at AT TIME ZONE 'Asia/Jakarta';

-- role
ALTER TABLE role
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- permission
ALTER TABLE permission
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta';

-- role_permission
ALTER TABLE role_permission
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta';

-- refresh_token (check if exists, table was dropped in migration 000049)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'refresh_token') THEN
        ALTER TABLE refresh_token
            ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
            ALTER COLUMN expired_at TYPE TIMESTAMP USING expired_at AT TIME ZONE 'Asia/Jakarta',
            ALTER COLUMN revoked_at TYPE TIMESTAMP USING revoked_at AT TIME ZONE 'Asia/Jakarta';
    END IF;
END $$;

-- activity_log
ALTER TABLE activity_log
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- 4. BUYER MODULE
-- =====================================================

-- buyer
ALTER TABLE buyer
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN last_login_at TYPE TIMESTAMP USING last_login_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN email_verified_at TYPE TIMESTAMP USING email_verified_at AT TIME ZONE 'Asia/Jakarta';

-- alamat_buyer
ALTER TABLE alamat_buyer
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- 5. WILAYAH (jika tabel masih ada)
-- =====================================================

-- Check if tables exist before altering
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'provinsi') THEN
        ALTER TABLE provinsi
            ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
            ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
            ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'kota') THEN
        ALTER TABLE kota
            ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
            ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
            ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'kecamatan') THEN
        ALTER TABLE kecamatan
            ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
            ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
            ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'kelurahan') THEN
        ALTER TABLE kelurahan
            ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
            ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
            ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';
    END IF;
END $$;

-- =====================================================
-- 6. MARKETING
-- =====================================================

-- hero_section
ALTER TABLE hero_section
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN tanggal_mulai TYPE TIMESTAMP USING tanggal_mulai AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN tanggal_selesai TYPE TIMESTAMP USING tanggal_selesai AT TIME ZONE 'Asia/Jakarta';

-- banner_event_promo
ALTER TABLE banner_event_promo
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN tanggal_mulai TYPE TIMESTAMP USING tanggal_mulai AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN tanggal_selesai TYPE TIMESTAMP USING tanggal_selesai AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- 7. OPERASIONAL
-- =====================================================

-- informasi_pickup (check if exists, table was dropped in migration 000095)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'informasi_pickup') THEN
        ALTER TABLE informasi_pickup
            ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
            ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta';
    END IF;
END $$;

-- jadwal_gudang
ALTER TABLE jadwal_gudang
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta';

-- dokumen_kebijakan
ALTER TABLE dokumen_kebijakan
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- disclaimer
ALTER TABLE disclaimer
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- 8. TRANSAKSI
-- =====================================================

-- metode_pembayaran_group
ALTER TABLE metode_pembayaran_group
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- metode_pembayaran
ALTER TABLE metode_pembayaran
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- ppn
ALTER TABLE ppn
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- pesanan
ALTER TABLE pesanan
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN expired_at TYPE TIMESTAMP USING expired_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN paid_at TYPE TIMESTAMP USING paid_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN processed_at TYPE TIMESTAMP USING processed_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN ready_at TYPE TIMESTAMP USING ready_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN shipped_at TYPE TIMESTAMP USING shipped_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN completed_at TYPE TIMESTAMP USING completed_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN cancelled_at TYPE TIMESTAMP USING cancelled_at AT TIME ZONE 'Asia/Jakarta';

-- pesanan_item
ALTER TABLE pesanan_item
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta';

-- pesanan_pembayaran
ALTER TABLE pesanan_pembayaran
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN expired_at TYPE TIMESTAMP USING expired_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN paid_at TYPE TIMESTAMP USING paid_at AT TIME ZONE 'Asia/Jakarta';

-- pesanan_status_history
ALTER TABLE pesanan_status_history
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- 9. ULASAN
-- =====================================================

-- ulasan
ALTER TABLE ulasan
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- 10. FORMULIR & KOMUNIKASI
-- =====================================================

-- formulir_partai_besar_config
ALTER TABLE formulir_partai_besar_config
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta';

-- formulir_partai_besar_anggaran
ALTER TABLE formulir_partai_besar_anggaran
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- formulir_partai_besar_submission
ALTER TABLE formulir_partai_besar_submission
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN email_sent_at TYPE TIMESTAMP USING email_sent_at AT TIME ZONE 'Asia/Jakarta';

-- whatsapp_handler
ALTER TABLE whatsapp_handler
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- 11. SISTEM KONTROL
-- =====================================================

-- force_update_app
ALTER TABLE force_update_app
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- mode_maintenance
ALTER TABLE mode_maintenance
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- 12. BLOG & VIDEO
-- =====================================================

-- kategori_blog
ALTER TABLE kategori_blog
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- label_blog
ALTER TABLE label_blog
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- blog
ALTER TABLE blog
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN published_at TYPE TIMESTAMP USING published_at AT TIME ZONE 'Asia/Jakarta';

-- kategori_video
ALTER TABLE kategori_video
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- video
ALTER TABLE video
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN published_at TYPE TIMESTAMP USING published_at AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- 13. FAQ
-- =====================================================

-- faq
ALTER TABLE faq
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'Asia/Jakarta',
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING deleted_at AT TIME ZONE 'Asia/Jakarta';

-- =====================================================
-- RECREATE ALL TRIGGERS
-- =====================================================
-- Recreate all triggers that were dropped at the beginning

-- Master Data
CREATE TRIGGER update_kategori_produk_updated_at
    BEFORE UPDATE ON kategori_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_merek_produk_updated_at
    BEFORE UPDATE ON merek_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_kondisi_produk_updated_at
    BEFORE UPDATE ON kondisi_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_kondisi_paket_updated_at
    BEFORE UPDATE ON kondisi_paket
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sumber_produk_updated_at
    BEFORE UPDATE ON sumber_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Produk Module
CREATE TRIGGER update_warehouse_updated_at
    BEFORE UPDATE ON warehouse
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tipe_produk_updated_at
    BEFORE UPDATE ON tipe_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_diskon_kategori_updated_at
    BEFORE UPDATE ON diskon_kategori
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_banner_tipe_produk_updated_at
    BEFORE UPDATE ON banner_tipe_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_produk_updated_at
    BEFORE UPDATE ON produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Admin & Auth
CREATE TRIGGER update_admin_updated_at
    BEFORE UPDATE ON admin
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_role_updated_at
    BEFORE UPDATE ON role
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_permission_updated_at
    BEFORE UPDATE ON permission
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Buyer Module
CREATE TRIGGER update_buyer_updated_at
    BEFORE UPDATE ON buyer
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_alamat_buyer_updated_at
    BEFORE UPDATE ON alamat_buyer
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_single_default_alamat
    AFTER INSERT OR UPDATE OF is_default ON alamat_buyer
    FOR EACH ROW
    WHEN (NEW.is_default = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_default_alamat();

CREATE TRIGGER trg_first_alamat_default
    BEFORE INSERT ON alamat_buyer
    FOR EACH ROW
    EXECUTE FUNCTION fn_first_alamat_as_default();

CREATE TRIGGER trg_prevent_delete_default
    BEFORE UPDATE ON alamat_buyer
    FOR EACH ROW
    EXECUTE FUNCTION fn_prevent_delete_default_alamat();

-- Wilayah (only if tables exist)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'provinsi') THEN
        EXECUTE 'CREATE TRIGGER update_provinsi_updated_at BEFORE UPDATE ON provinsi FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'kota') THEN
        EXECUTE 'CREATE TRIGGER update_kota_updated_at BEFORE UPDATE ON kota FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'kecamatan') THEN
        EXECUTE 'CREATE TRIGGER update_kecamatan_updated_at BEFORE UPDATE ON kecamatan FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'kelurahan') THEN
        EXECUTE 'CREATE TRIGGER update_kelurahan_updated_at BEFORE UPDATE ON kelurahan FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()';
    END IF;
END $$;

-- Marketing
CREATE TRIGGER trg_hero_section_updated_at
    BEFORE UPDATE ON hero_section
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_banner_event_updated_at
    BEFORE UPDATE ON banner_event_promo
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Transaksi
CREATE TRIGGER trg_metode_pembayaran_group_updated_at
    BEFORE UPDATE ON metode_pembayaran_group
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_metode_pembayaran_updated_at
    BEFORE UPDATE ON metode_pembayaran
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_ppn_updated_at
    BEFORE UPDATE ON ppn
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_pesanan_updated_at
    BEFORE UPDATE ON pesanan
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_pesanan_item_updated_at
    BEFORE UPDATE ON pesanan_item
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_pesanan_pembayaran_updated_at
    BEFORE UPDATE ON pesanan_pembayaran
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Operasional
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'informasi_pickup') THEN
        EXECUTE 'CREATE TRIGGER trg_informasi_pickup_updated_at BEFORE UPDATE ON informasi_pickup FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()';
    END IF;
END $$;

CREATE TRIGGER trg_jadwal_gudang_updated_at
    BEFORE UPDATE ON jadwal_gudang
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_dokumen_kebijakan_updated_at
    BEFORE UPDATE ON dokumen_kebijakan
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_disclaimer_updated_at
    BEFORE UPDATE ON disclaimer
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Ulasan
CREATE TRIGGER trg_ulasan_updated_at
    BEFORE UPDATE ON ulasan
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Formulir & Komunikasi
CREATE TRIGGER trg_formulir_config_updated_at
    BEFORE UPDATE ON formulir_partai_besar_config
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_formulir_anggaran_updated_at
    BEFORE UPDATE ON formulir_partai_besar_anggaran
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_whatsapp_handler_updated_at
    BEFORE UPDATE ON whatsapp_handler
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Sistem Kontrol
CREATE TRIGGER trg_force_update_updated_at
    BEFORE UPDATE ON force_update_app
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_maintenance_updated_at
    BEFORE UPDATE ON mode_maintenance
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Blog & Video
CREATE TRIGGER trg_kategori_blog_updated_at
    BEFORE UPDATE ON kategori_blog
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_label_blog_updated_at
    BEFORE UPDATE ON label_blog
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_blog_updated_at
    BEFORE UPDATE ON blog
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_kategori_video_updated_at
    BEFORE UPDATE ON kategori_video
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_video_updated_at
    BEFORE UPDATE ON video
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- FAQ
CREATE TRIGGER trg_faq_updated_at
    BEFORE UPDATE ON faq
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- RECREATE OTHER BUSINESS LOGIC TRIGGERS
-- =====================================================
-- Triggers that reference deleted_at or timestamp columns

-- Disclaimer triggers
CREATE TRIGGER trg_generate_slug_disclaimer
    BEFORE INSERT OR UPDATE ON disclaimer
    FOR EACH ROW 
    EXECUTE FUNCTION generate_slug_disclaimer();

CREATE TRIGGER trg_single_active_disclaimer
    AFTER INSERT OR UPDATE OF is_active ON disclaimer
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_disclaimer();

CREATE TRIGGER trigger_rewrite_slug_on_delete_disclaimer
    BEFORE UPDATE ON disclaimer
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- WhatsApp Handler
CREATE TRIGGER trg_single_active_wa_handler
    AFTER INSERT OR UPDATE OF is_active ON whatsapp_handler
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_wa_handler();

-- Mode Maintenance
CREATE TRIGGER trg_single_active_maintenance
    AFTER INSERT OR UPDATE OF is_active ON mode_maintenance
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_maintenance();

-- Force Update App
CREATE TRIGGER trg_single_active_force_update
    AFTER INSERT OR UPDATE OF is_active ON force_update_app
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_force_update();

-- PPN
CREATE TRIGGER trg_single_active_ppn
    AFTER INSERT OR UPDATE OF is_active ON ppn
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_ppn();

-- Hero Section
CREATE TRIGGER trg_hero_section_auto_sync_insert
    BEFORE INSERT ON hero_section
    FOR EACH ROW
    EXECUTE FUNCTION fn_hero_section_auto_sync_insert();

-- Formulir Config
CREATE TRIGGER trg_formulir_config_singleton
    BEFORE INSERT ON formulir_partai_besar_config
    FOR EACH ROW
    EXECUTE FUNCTION check_formulir_config_singleton();

-- Informasi Pickup (only if exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'informasi_pickup') THEN
        EXECUTE 'CREATE TRIGGER trg_informasi_pickup_singleton BEFORE INSERT ON informasi_pickup FOR EACH ROW EXECUTE FUNCTION check_informasi_pickup_singleton()';
    END IF;
END $$;

-- Pesanan
CREATE TRIGGER trg_log_pesanan_status_change
    AFTER UPDATE ON pesanan
    FOR EACH ROW
    EXECUTE FUNCTION log_pesanan_status_change();

-- Pesanan Item
CREATE TRIGGER trg_calculate_subtotal
    BEFORE INSERT OR UPDATE ON pesanan_item
    FOR EACH ROW
    EXECUTE FUNCTION calculate_pesanan_item_subtotal();

-- Pesanan Pembayaran
CREATE TRIGGER trg_check_max_split_payment
    BEFORE INSERT ON pesanan_pembayaran
    FOR EACH ROW
    EXECUTE FUNCTION check_max_split_payment();

-- Ulasan
CREATE TRIGGER trg_validate_ulasan_order
    BEFORE INSERT ON ulasan
    FOR EACH ROW
    EXECUTE FUNCTION fn_validate_ulasan_order_completed();

CREATE TRIGGER trg_set_ulasan_approved_at
    BEFORE UPDATE ON ulasan
    FOR EACH ROW
    EXECUTE FUNCTION fn_set_ulasan_approved_at();

-- Blog
CREATE TRIGGER trg_blog_published_at
    BEFORE UPDATE ON blog
    FOR EACH ROW
    EXECUTE FUNCTION set_blog_published_at();

-- Video
CREATE TRIGGER trg_video_published_at
    BEFORE UPDATE ON video
    FOR EACH ROW
    EXECUTE FUNCTION set_video_published_at();

-- Rewrite slug on delete triggers for master data
CREATE TRIGGER trigger_rewrite_slug_on_delete_kategori_produk
    BEFORE UPDATE ON kategori_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

CREATE TRIGGER trigger_rewrite_slug_on_delete_merek_produk
    BEFORE UPDATE ON merek_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

CREATE TRIGGER trigger_rewrite_slug_on_delete_kondisi_produk
    BEFORE UPDATE ON kondisi_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

CREATE TRIGGER trigger_rewrite_slug_on_delete_kondisi_paket
    BEFORE UPDATE ON kondisi_paket
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

CREATE TRIGGER trigger_rewrite_slug_on_delete_sumber_produk
    BEFORE UPDATE ON sumber_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

CREATE TRIGGER trigger_rewrite_slug_on_delete_warehouse
    BEFORE UPDATE ON warehouse
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

CREATE TRIGGER trigger_rewrite_slug_on_delete_tipe_produk
    BEFORE UPDATE ON tipe_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

CREATE TRIGGER trigger_rewrite_slug_on_delete_produk
    BEFORE UPDATE ON produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- =====================================================
-- VERIFICATION
-- =====================================================

-- Log completion
DO $$
BEGIN
    RAISE NOTICE 'Rollback completed: All TIMESTAMPTZ columns converted back to TIMESTAMP';
    RAISE NOTICE 'Timezone: Asia/Jakarta (UTC+7)';
    RAISE NOTICE 'WARNING: Timezone information has been lost';
END $$;
