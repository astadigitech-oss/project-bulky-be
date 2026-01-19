-- =====================================================
-- SEED: WhatsApp Handler
-- =====================================================
-- Insert default WhatsApp handler configuration
-- Primary: 62811833164 (Active)
-- Backup: 6289876543210 (Inactive)
-- =====================================================

-- Add unique constraint if not exists
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'whatsapp_handler_nomor_wa_key'
    ) THEN
        ALTER TABLE whatsapp_handler 
        ADD CONSTRAINT whatsapp_handler_nomor_wa_key UNIQUE (nomor_wa);
    END IF;
END $$;

-- Insert primary WhatsApp handler (active)
INSERT INTO whatsapp_handler (id, nomor_wa, pesan_awal, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    '62811833164',
    'Halo, ada yang bisa kami bantu?',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (nomor_wa) DO UPDATE SET
    pesan_awal = EXCLUDED.pesan_awal,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Insert backup WhatsApp handler (inactive)
INSERT INTO whatsapp_handler (id, nomor_wa, pesan_awal, is_active, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    '6289876543210',
    'Halo, tim Bulky siap membantu Anda!',
    false,
    NOW(),
    NOW()
)
ON CONFLICT (nomor_wa) DO UPDATE SET
    pesan_awal = EXCLUDED.pesan_awal,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();

-- Add comment
COMMENT ON TABLE whatsapp_handler IS 'WhatsApp handler configuration for customer support';
