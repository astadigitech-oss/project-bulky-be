-- =====================================================
-- ALTER: Convert nama JSON to nama_id and nama_en VARCHAR
-- =====================================================

-- Add new columns if not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'kategori_video' AND column_name = 'nama_id') THEN
        ALTER TABLE kategori_video ADD COLUMN nama_id VARCHAR(100);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'kategori_video' AND column_name = 'nama_en') THEN
        ALTER TABLE kategori_video ADD COLUMN nama_en VARCHAR(100);
    END IF;
END $$;

-- Migrate data from JSON field 'nama' to separate columns (if nama field exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'kategori_video' AND column_name = 'nama') THEN
        
        -- Copy data from JSON to new fields
        UPDATE kategori_video 
        SET nama_id = nama->>'id',
            nama_en = nama->>'en'
        WHERE nama IS NOT NULL;
        
        -- Drop old JSON field
        ALTER TABLE kategori_video DROP COLUMN nama;
    END IF;
END $$;

-- Make nama_id NOT NULL after data migration
ALTER TABLE kategori_video ALTER COLUMN nama_id SET NOT NULL;
