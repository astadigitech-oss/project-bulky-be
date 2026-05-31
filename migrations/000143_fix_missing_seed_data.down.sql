-- Reverse: reset bilingual content to ID content and clear slug_en/coordinates

UPDATE dokumen_kebijakan SET judul_en = judul, konten_en = konten WHERE deleted_at IS NULL;
UPDATE dokumen_kebijakan SET slug_en = NULL WHERE deleted_at IS NULL;
UPDATE warehouse SET latitude = NULL, longitude = NULL WHERE slug = 'warehouse-cibinong';
UPDATE disclaimer SET slug_en = NULL WHERE slug_id = 'disclaimer-pembelian';
