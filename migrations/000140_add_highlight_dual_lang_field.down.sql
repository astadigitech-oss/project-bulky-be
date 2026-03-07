DROP INDEX IF EXISTS idx_blog_highlight_en;
DROP INDEX IF EXISTS idx_blog_highlight_id;

ALTER TABLE blog DROP COLUMN IF EXISTS highlight_en;
ALTER TABLE blog DROP COLUMN IF EXISTS highlight_id;