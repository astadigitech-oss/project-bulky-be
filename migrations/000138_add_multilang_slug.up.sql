-- Add multi-lang slug columns (slug_id, slug_en) to 13 tables
-- slug_id: Indonesian slug, slug_en: English slug
-- Both nullable, each uniquely indexed
-- Existing slug data migrated to slug_id

-- ============================================
-- 1. blog
-- ============================================
ALTER TABLE blog ADD COLUMN IF NOT EXISTS slug_id VARCHAR(250);
ALTER TABLE blog ADD COLUMN IF NOT EXISTS slug_en VARCHAR(250);

UPDATE blog SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_blog_slug_id ON blog(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_blog_slug_en ON blog(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 2. disclaimer
-- ============================================
ALTER TABLE disclaimer ADD COLUMN IF NOT EXISTS slug_id VARCHAR(200);
ALTER TABLE disclaimer ADD COLUMN IF NOT EXISTS slug_en VARCHAR(200);

UPDATE disclaimer SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_disclaimer_slug_id ON disclaimer(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_disclaimer_slug_en ON disclaimer(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 3. dokumen_kebijakan
-- ============================================
ALTER TABLE dokumen_kebijakan ADD COLUMN IF NOT EXISTS slug_id VARCHAR(100);
ALTER TABLE dokumen_kebijakan ADD COLUMN IF NOT EXISTS slug_en VARCHAR(100);

UPDATE dokumen_kebijakan SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_dokumen_kebijakan_slug_id ON dokumen_kebijakan(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_dokumen_kebijakan_slug_en ON dokumen_kebijakan(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 4. kategori_produk
-- ============================================
ALTER TABLE kategori_produk ADD COLUMN IF NOT EXISTS slug_id VARCHAR(120);
ALTER TABLE kategori_produk ADD COLUMN IF NOT EXISTS slug_en VARCHAR(120);

UPDATE kategori_produk SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_kategori_produk_slug_id ON kategori_produk(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_kategori_produk_slug_en ON kategori_produk(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 5. kategori_blog
-- ============================================
ALTER TABLE kategori_blog ADD COLUMN IF NOT EXISTS slug_id VARCHAR(100);
ALTER TABLE kategori_blog ADD COLUMN IF NOT EXISTS slug_en VARCHAR(100);

UPDATE kategori_blog SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_kategori_blog_slug_id ON kategori_blog(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_kategori_blog_slug_en ON kategori_blog(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 6. kategori_video
-- ============================================
ALTER TABLE kategori_video ADD COLUMN IF NOT EXISTS slug_id VARCHAR(100);
ALTER TABLE kategori_video ADD COLUMN IF NOT EXISTS slug_en VARCHAR(100);

UPDATE kategori_video SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_kategori_video_slug_id ON kategori_video(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_kategori_video_slug_en ON kategori_video(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 7. kondisi_produk
-- ============================================
ALTER TABLE kondisi_produk ADD COLUMN IF NOT EXISTS slug_id VARCHAR(120);
ALTER TABLE kondisi_produk ADD COLUMN IF NOT EXISTS slug_en VARCHAR(120);

UPDATE kondisi_produk SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_kondisi_produk_slug_id ON kondisi_produk(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_kondisi_produk_slug_en ON kondisi_produk(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 8. kondisi_paket
-- ============================================
ALTER TABLE kondisi_paket ADD COLUMN IF NOT EXISTS slug_id VARCHAR(120);
ALTER TABLE kondisi_paket ADD COLUMN IF NOT EXISTS slug_en VARCHAR(120);

UPDATE kondisi_paket SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_kondisi_paket_slug_id ON kondisi_paket(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_kondisi_paket_slug_en ON kondisi_paket(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 9. label_blog
-- ============================================
ALTER TABLE label_blog ADD COLUMN IF NOT EXISTS slug_id VARCHAR(100);
ALTER TABLE label_blog ADD COLUMN IF NOT EXISTS slug_en VARCHAR(100);

UPDATE label_blog SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_label_blog_slug_id ON label_blog(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_label_blog_slug_en ON label_blog(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 10. merek_produk
-- ============================================
ALTER TABLE merek_produk ADD COLUMN IF NOT EXISTS slug_id VARCHAR(120);
ALTER TABLE merek_produk ADD COLUMN IF NOT EXISTS slug_en VARCHAR(120);

UPDATE merek_produk SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_merek_produk_slug_id ON merek_produk(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_merek_produk_slug_en ON merek_produk(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 11. produk
-- ============================================
ALTER TABLE produk ADD COLUMN IF NOT EXISTS slug_id VARCHAR(280);
ALTER TABLE produk ADD COLUMN IF NOT EXISTS slug_en VARCHAR(280);

UPDATE produk SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_produk_slug_id ON produk(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_produk_slug_en ON produk(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 12. sumber_produk
-- ============================================
ALTER TABLE sumber_produk ADD COLUMN IF NOT EXISTS slug_id VARCHAR(120);
ALTER TABLE sumber_produk ADD COLUMN IF NOT EXISTS slug_en VARCHAR(120);

UPDATE sumber_produk SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_sumber_produk_slug_id ON sumber_produk(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_sumber_produk_slug_en ON sumber_produk(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;

-- ============================================
-- 13. video
-- ============================================
ALTER TABLE video ADD COLUMN IF NOT EXISTS slug_id VARCHAR(250);
ALTER TABLE video ADD COLUMN IF NOT EXISTS slug_en VARCHAR(250);

UPDATE video SET slug_id = slug WHERE slug IS NOT NULL AND slug_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_video_slug_id ON video(slug_id) WHERE slug_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_video_slug_en ON video(slug_en) WHERE slug_en IS NOT NULL AND deleted_at IS NULL;
