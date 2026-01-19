-- =====================================================
-- FIX: Slug Unique Constraint on Soft Delete
-- =====================================================
-- Saat soft delete, append timestamp ke slug
-- Format: {slug}-deleted-{unix_timestamp_microseconds}
-- Menggunakan microseconds untuk menghindari collision
-- =====================================================

-- Function: Rewrite slug on soft delete
CREATE OR REPLACE FUNCTION rewrite_slug_on_delete()
RETURNS TRIGGER AS $$
DECLARE
    new_slug TEXT;
    counter INTEGER := 0;
BEGIN
    -- Cek apakah deleted_at berubah dari NULL ke NOT NULL
    IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
        -- Generate slug dengan timestamp microseconds
        new_slug := OLD.slug || '-deleted-' || 
                    EXTRACT(EPOCH FROM NEW.deleted_at)::bigint || 
                    LPAD(EXTRACT(MICROSECONDS FROM NEW.deleted_at)::TEXT, 6, '0');
        
        -- Jika masih ada collision (sangat jarang), tambahkan counter
        -- Loop sampai dapat slug yang unique
        WHILE EXISTS (
            SELECT 1 FROM (
                SELECT slug FROM kategori_produk WHERE slug = new_slug
                UNION ALL
                SELECT slug FROM merek_produk WHERE slug = new_slug
                UNION ALL
                SELECT slug FROM kondisi_produk WHERE slug = new_slug
                UNION ALL
                SELECT slug FROM kondisi_paket WHERE slug = new_slug
                UNION ALL
                SELECT slug FROM sumber_produk WHERE slug = new_slug
                UNION ALL
                SELECT slug FROM warehouse WHERE slug = new_slug
                UNION ALL
                SELECT slug FROM tipe_produk WHERE slug = new_slug
                UNION ALL
                SELECT slug FROM produk WHERE slug = new_slug
                UNION ALL
                SELECT slug FROM dokumen_kebijakan WHERE slug = new_slug
                UNION ALL
                SELECT slug FROM disclaimer WHERE slug = new_slug
            ) AS all_slugs
            LIMIT 1
        ) LOOP
            counter := counter + 1;
            new_slug := OLD.slug || '-deleted-' || 
                       EXTRACT(EPOCH FROM NEW.deleted_at)::bigint || 
                       LPAD(EXTRACT(MICROSECONDS FROM NEW.deleted_at)::TEXT, 6, '0') ||
                       '-' || counter;
        END LOOP;
        
        NEW.slug := new_slug;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Apply trigger ke semua tabel dengan slug
-- =====================================================

-- kategori_produk
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_kategori_produk ON kategori_produk;
CREATE TRIGGER trigger_rewrite_slug_on_delete_kategori_produk
    BEFORE UPDATE ON kategori_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- merek_produk
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_merek_produk ON merek_produk;
CREATE TRIGGER trigger_rewrite_slug_on_delete_merek_produk
    BEFORE UPDATE ON merek_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- kondisi_produk
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_kondisi_produk ON kondisi_produk;
CREATE TRIGGER trigger_rewrite_slug_on_delete_kondisi_produk
    BEFORE UPDATE ON kondisi_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- kondisi_paket
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_kondisi_paket ON kondisi_paket;
CREATE TRIGGER trigger_rewrite_slug_on_delete_kondisi_paket
    BEFORE UPDATE ON kondisi_paket
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- sumber_produk
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_sumber_produk ON sumber_produk;
CREATE TRIGGER trigger_rewrite_slug_on_delete_sumber_produk
    BEFORE UPDATE ON sumber_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- warehouse
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_warehouse ON warehouse;
CREATE TRIGGER trigger_rewrite_slug_on_delete_warehouse
    BEFORE UPDATE ON warehouse
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- tipe_produk
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_tipe_produk ON tipe_produk;
CREATE TRIGGER trigger_rewrite_slug_on_delete_tipe_produk
    BEFORE UPDATE ON tipe_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- produk
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_produk ON produk;
CREATE TRIGGER trigger_rewrite_slug_on_delete_produk
    BEFORE UPDATE ON produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- dokumen_kebijakan
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_dokumen_kebijakan ON dokumen_kebijakan;
CREATE TRIGGER trigger_rewrite_slug_on_delete_dokumen_kebijakan
    BEFORE UPDATE ON dokumen_kebijakan
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- disclaimer
DROP TRIGGER IF EXISTS trigger_rewrite_slug_on_delete_disclaimer ON disclaimer;
CREATE TRIGGER trigger_rewrite_slug_on_delete_disclaimer
    BEFORE UPDATE ON disclaimer
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- =====================================================
-- Add comment
-- =====================================================
COMMENT ON FUNCTION rewrite_slug_on_delete() IS 
'Automatically append timestamp with microseconds to slug when soft deleting to prevent unique constraint violation. Includes collision detection with counter fallback.';
