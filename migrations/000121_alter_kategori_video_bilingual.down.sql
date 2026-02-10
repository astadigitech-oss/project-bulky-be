-- Rollback
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'kategori_video' AND column_name = 'nama_en') THEN
        ALTER TABLE kategori_video DROP COLUMN IF EXISTS nama_en;
    END IF;
END $$;
