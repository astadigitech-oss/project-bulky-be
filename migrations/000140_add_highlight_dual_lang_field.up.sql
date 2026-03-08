-- Add multi-lang highlight columns (highlight_id, highlight_en) to table blog
-- highlight_id: Indonesian highlight, highlight_en: English highlight
-- Both nullable, each uniquely indexed

ALTER TABLE blog ADD COLUMN IF NOT EXISTS highlight_id TEXT;
ALTER TABLE blog ADD COLUMN IF NOT EXISTS highlight_en TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_blog_highlight_id ON blog(highlight_id) WHERE highlight_id IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_blog_highlight_en ON blog(highlight_en) WHERE highlight_en IS NOT NULL AND deleted_at IS NULL;