-- Restore nomor_whatsapp column to warehouse
ALTER TABLE warehouse
ADD COLUMN IF NOT EXISTS nomor_whatsapp VARCHAR(20);

-- Copy telepon to nomor_whatsapp
UPDATE warehouse
SET nomor_whatsapp = telepon
WHERE telepon IS NOT NULL;

-- Add comment
COMMENT ON COLUMN warehouse.nomor_whatsapp IS 'Nomor WhatsApp untuk kontak (format: 62xxx)';
