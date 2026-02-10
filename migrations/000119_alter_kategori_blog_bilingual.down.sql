-- Rollback: Remove added columns if they didn't exist before
-- Note: This won't restore renamed columns

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'kategori_blog' AND column_name = 'nama_en') THEN
        ALTER TABLE kategori_blog DROP COLUMN IF EXISTS nama_en;
    END IF;
END $$;
