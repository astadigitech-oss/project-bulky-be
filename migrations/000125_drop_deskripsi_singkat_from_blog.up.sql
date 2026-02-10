-- Drop the search index that includes deskripsi_singkat_id
DROP INDEX IF EXISTS idx_blog_search;

-- Drop the columns
ALTER TABLE blog 
    DROP COLUMN IF EXISTS deskripsi_singkat_id,
    DROP COLUMN IF EXISTS deskripsi_singkat_en;

-- Recreate search index without deskripsi_singkat_id
CREATE INDEX idx_blog_search ON blog USING gin(to_tsvector('indonesian', judul_id));
