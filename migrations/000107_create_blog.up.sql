CREATE TABLE blog (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    judul_id VARCHAR(200) NOT NULL,
    judul_en VARCHAR(200),
    slug VARCHAR(250) UNIQUE NOT NULL,
    konten_id TEXT NOT NULL,
    konten_en TEXT,
    deskripsi_singkat_id TEXT NOT NULL,
    deskripsi_singkat_en TEXT,
    featured_image_url VARCHAR(500),
    kategori_id UUID NOT NULL,
    meta_title_id VARCHAR(200),
    meta_title_en VARCHAR(200),
    meta_description_id TEXT,
    meta_description_en TEXT,
    meta_keywords TEXT,
    is_active BOOLEAN DEFAULT false,
    view_count INTEGER DEFAULT 0,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT fk_blog_kategori FOREIGN KEY (kategori_id) REFERENCES kategori_blog(id) ON DELETE RESTRICT,
    CONSTRAINT chk_blog_konten_not_empty CHECK (LENGTH(TRIM(konten_id)) > 0)
);

CREATE INDEX idx_blog_slug ON blog(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_blog_kategori ON blog(kategori_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_blog_active ON blog(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_blog_published ON blog(published_at DESC) WHERE deleted_at IS NULL AND is_active = true;
CREATE INDEX idx_blog_view_count ON blog(view_count DESC);
CREATE INDEX idx_blog_search ON blog USING gin(to_tsvector('indonesian', judul_id || ' ' || deskripsi_singkat_id));

CREATE TRIGGER trg_blog_updated_at
  BEFORE UPDATE ON blog
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- Auto set published_at
CREATE OR REPLACE FUNCTION set_blog_published_at()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.is_active = true AND OLD.is_active = false AND NEW.published_at IS NULL THEN
    NEW.published_at = NOW();
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_blog_published_at
  BEFORE UPDATE ON blog
  FOR EACH ROW
  EXECUTE FUNCTION set_blog_published_at();
