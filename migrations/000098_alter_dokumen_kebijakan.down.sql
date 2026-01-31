-- Add urutan column back (only if not exists)
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'dokumen_kebijakan' AND column_name = 'urutan'
    ) THEN
        ALTER TABLE dokumen_kebijakan ADD COLUMN urutan INT DEFAULT 0;
    END IF;
END $$;

-- Recreate index
DROP INDEX IF EXISTS idx_dokumen_kebijakan_urutan;
CREATE INDEX idx_dokumen_kebijakan_urutan ON dokumen_kebijakan(urutan) 
    WHERE deleted_at IS NULL;
