-- Drop nomor_whatsapp column from warehouse (use telepon instead)
ALTER TABLE warehouse
DROP COLUMN IF EXISTS nomor_whatsapp;

-- Update comment on telepon column
COMMENT ON COLUMN warehouse.telepon IS 'Nomor telepon/WhatsApp untuk kontak (format: 62xxx)';
