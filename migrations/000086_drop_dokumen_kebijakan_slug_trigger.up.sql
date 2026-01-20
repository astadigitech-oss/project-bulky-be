-- =====================================================
-- DROP: Dokumen Kebijakan Slug Trigger
-- =====================================================
-- Remove auto-slug generation trigger since dokumen_kebijakan
-- no longer uses slug column (removed in migration 000085)
-- =====================================================

-- Drop the trigger
DROP TRIGGER IF EXISTS trg_generate_slug_dokumen ON dokumen_kebijakan;

-- Drop the function
DROP FUNCTION IF EXISTS generate_slug_dokumen_kebijakan();
