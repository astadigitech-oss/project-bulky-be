-- Rollback
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'label_blog' AND column_name = 'nama_en') THEN
        ALTER TABLE label_blog DROP COLUMN IF EXISTS nama_en;
    END IF;
END $$;
