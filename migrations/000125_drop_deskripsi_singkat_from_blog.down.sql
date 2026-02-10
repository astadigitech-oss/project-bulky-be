-- Drop the new search index
DROP INDEX IF EXISTS idx_blog_search;

-- Add back the columns
ALTER TABLE blog 
    ADD COLUMN deskripsi_singkat_id TEXT NOT NULL DEFAULT '',
    ADD COLUMN deskripsi_singkat_en TEXT;

-- Recreate search index with deskripsi_singkat_id
CREATE INDEX idx_blog_search ON blog USING gin(to_tsvector('indonesian', judul_id || ' ' || deskripsi_singkat_id));
