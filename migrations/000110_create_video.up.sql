CREATE TABLE video (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    judul_id VARCHAR(200) NOT NULL,
    judul_en VARCHAR(200),
    slug VARCHAR(250) UNIQUE NOT NULL,
    deskripsi_id TEXT NOT NULL,
    deskripsi_en TEXT,
    video_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    kategori_id UUID NOT NULL,
    durasi_detik INTEGER NOT NULL,
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
    CONSTRAINT fk_video_kategori FOREIGN KEY (kategori_id) REFERENCES kategori_video(id) ON DELETE RESTRICT,
    CONSTRAINT chk_video_durasi_positive CHECK (durasi_detik > 0)
);

CREATE INDEX idx_video_slug ON video(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_video_kategori ON video(kategori_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_video_active ON video(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_video_published ON video(published_at DESC) WHERE deleted_at IS NULL AND is_active = true;
CREATE INDEX idx_video_view_count ON video(view_count DESC);
CREATE INDEX idx_video_search ON video USING gin(to_tsvector('indonesian', judul_id || ' ' || deskripsi_id));

CREATE TRIGGER trg_video_updated_at
  BEFORE UPDATE ON video
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- Auto set published_at
CREATE OR REPLACE FUNCTION set_video_published_at()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.is_active = true AND OLD.is_active = false AND NEW.published_at IS NULL THEN
    NEW.published_at = NOW();
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_video_published_at
  BEFORE UPDATE ON video
  FOR EACH ROW
  EXECUTE FUNCTION set_video_published_at();
