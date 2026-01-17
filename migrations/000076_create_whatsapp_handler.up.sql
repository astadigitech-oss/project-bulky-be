-- =====================================================
-- TABEL: whatsapp_handler
-- =====================================================
-- Konfigurasi WhatsApp floating button
-- Hanya 1 yang bisa aktif
-- =====================================================

CREATE TABLE whatsapp_handler (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nomor_wa VARCHAR(20) NOT NULL,
    pesan_awal TEXT NOT NULL,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_whatsapp_handler_is_active ON whatsapp_handler(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_whatsapp_handler_deleted_at ON whatsapp_handler(deleted_at);

CREATE TRIGGER trg_whatsapp_handler_updated_at
    BEFORE UPDATE ON whatsapp_handler
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger: Ensure single active
CREATE OR REPLACE FUNCTION fn_ensure_single_active_wa_handler()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = true AND NEW.deleted_at IS NULL THEN
        UPDATE whatsapp_handler 
        SET is_active = false, updated_at = NOW()
        WHERE id != NEW.id AND is_active = true AND deleted_at IS NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_single_active_wa_handler
    AFTER INSERT OR UPDATE OF is_active ON whatsapp_handler
    FOR EACH ROW
    WHEN (NEW.is_active = true AND NEW.deleted_at IS NULL)
    EXECUTE FUNCTION fn_ensure_single_active_wa_handler();

COMMENT ON TABLE whatsapp_handler IS 'Konfigurasi WhatsApp floating button. Hanya 1 yang bisa aktif.';
COMMENT ON COLUMN whatsapp_handler.nomor_wa IS 'Nomor WhatsApp, format 62xxx (tanpa +)';
COMMENT ON COLUMN whatsapp_handler.pesan_awal IS 'Template pesan yang muncul saat buyer klik WA button';
