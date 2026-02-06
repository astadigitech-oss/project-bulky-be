CREATE TABLE label_blog (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama_id VARCHAR(100) NOT NULL,
    nama_en VARCHAR(100),
    slug VARCHAR(100) UNIQUE NOT NULL,
    urutan INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_label_blog_slug ON label_blog(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_label_blog_urutan ON label_blog(urutan);

CREATE TRIGGER trg_label_blog_updated_at
  BEFORE UPDATE ON label_blog
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

