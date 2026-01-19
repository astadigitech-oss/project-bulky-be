DROP TRIGGER IF EXISTS trg_single_active_wa_handler ON whatsapp_handler;
DROP TRIGGER IF EXISTS trg_whatsapp_handler_updated_at ON whatsapp_handler;
DROP FUNCTION IF EXISTS fn_ensure_single_active_wa_handler();
DROP TABLE IF EXISTS whatsapp_handler;
