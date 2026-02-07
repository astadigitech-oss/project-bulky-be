CREATE TABLE kategori_blog (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama_id VARCHAR(100) NOT NULL,
    nama_en VARCHAR(100),
    slug VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    urutan INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_kategori_blog_slug ON kategori_blog(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_kategori_blog_active ON kategori_blog(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_kategori_blog_urutan ON kategori_blog(urutan);

CREATE TRIGGER trg_kategori_blog_updated_at
  BEFORE UPDATE ON kategori_blog
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

