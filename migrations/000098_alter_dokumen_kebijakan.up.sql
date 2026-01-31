-- Hapus field urutan (only if exists)
DO $$ 
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'dokumen_kebijakan' AND column_name = 'urutan'
    ) THEN
        ALTER TABLE dokumen_kebijakan DROP COLUMN urutan;
    END IF;
END $$;

-- Drop related index if exists
DROP INDEX IF EXISTS idx_dokumen_kebijakan_urutan;
