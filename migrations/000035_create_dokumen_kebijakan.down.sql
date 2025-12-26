DROP TRIGGER IF EXISTS trg_generate_slug_dokumen ON dokumen_kebijakan;
DROP TRIGGER IF EXISTS trg_dokumen_kebijakan_updated_at ON dokumen_kebijakan;
DROP FUNCTION IF EXISTS generate_slug_dokumen_kebijakan();
DROP TABLE IF EXISTS dokumen_kebijakan;
