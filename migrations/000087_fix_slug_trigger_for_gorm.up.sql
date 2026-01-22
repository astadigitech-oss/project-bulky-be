-- =====================================================
-- FIX: Slug Trigger untuk GORM Soft Delete
-- =====================================================
-- GORM soft delete hanya update deleted_at tanpa load kolom lain
-- Kita perlu query slug dari tabel saat trigger dijalankan
-- =====================================================

CREATE OR REPLACE FUNCTION rewrite_slug_on_delete()
RETURNS TRIGGER AS $$
DECLARE
    new_slug TEXT;
    old_slug TEXT;
    counter INTEGER := 0;
BEGIN
    -- Cek apakah deleted_at berubah dari NULL ke NOT NULL (soft delete)
    IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
        -- Ambil slug dari NEW jika ada, fallback ke query tabel
        -- Ini handle case dimana GORM tidak pass semua kolom
        IF NEW.slug IS NOT NULL AND NEW.slug != '' THEN
            old_slug := NEW.slug;
        ELSE
            -- Query slug dari tabel berdasarkan primary key
            EXECUTE format('SELECT slug FROM %I WHERE id = $1', TG_TABLE_NAME)
            INTO old_slug
            USING NEW.id;
        END IF;
        
        -- Jika slug tidak ditemukan, skip rewrite
        IF old_slug IS NULL OR old_slug = '' THEN
            RETURN NEW;
        END IF;
        
        -- Generate slug dengan timestamp microseconds
        new_slug := old_slug || '-deleted-' || 
                    EXTRACT(EPOCH FROM NEW.deleted_at)::bigint || 
                    LPAD(EXTRACT(MICROSECONDS FROM NEW.deleted_at)::TEXT, 6, '0');
        
        -- Jika masih ada collision (sangat jarang), tambahkan counter
        WHILE EXISTS (
            SELECT 1 FROM (
                SELECT slug FROM kategori_produk WHERE slug = new_slug AND deleted_at IS NULL
                UNION ALL
                SELECT slug FROM merek_produk WHERE slug = new_slug AND deleted_at IS NULL
                UNION ALL
                SELECT slug FROM kondisi_produk WHERE slug = new_slug AND deleted_at IS NULL
                UNION ALL
                SELECT slug FROM kondisi_paket WHERE slug = new_slug AND deleted_at IS NULL
                UNION ALL
                SELECT slug FROM sumber_produk WHERE slug = new_slug AND deleted_at IS NULL
                UNION ALL
                SELECT slug FROM warehouse WHERE slug = new_slug AND deleted_at IS NULL
                UNION ALL
                SELECT slug FROM tipe_produk WHERE slug = new_slug AND deleted_at IS NULL
                UNION ALL
                SELECT slug FROM produk WHERE slug = new_slug AND deleted_at IS NULL
                UNION ALL
                SELECT slug FROM dokumen_kebijakan WHERE slug = new_slug AND deleted_at IS NULL
                UNION ALL
                SELECT slug FROM disclaimer WHERE slug = new_slug AND deleted_at IS NULL
            ) AS all_slugs
            LIMIT 1
        ) LOOP
            counter := counter + 1;
            new_slug := old_slug || '-deleted-' || 
                       EXTRACT(EPOCH FROM NEW.deleted_at)::bigint || 
                       LPAD(EXTRACT(MICROSECONDS FROM NEW.deleted_at)::TEXT, 6, '0') ||
                       '-' || counter;
        END LOOP;
        
        NEW.slug := new_slug;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION rewrite_slug_on_delete() IS 
'Automatically append timestamp with microseconds to slug when soft deleting. Handles GORM soft delete that only updates deleted_at column.';
