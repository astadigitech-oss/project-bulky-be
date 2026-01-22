-- Rollback ke versi trigger sebelumnya
CREATE OR REPLACE FUNCTION rewrite_slug_on_delete()
RETURNS TRIGGER AS $
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
$ LANGUAGE plpgsql;
