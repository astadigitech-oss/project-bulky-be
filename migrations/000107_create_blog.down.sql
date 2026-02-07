DROP TRIGGER IF EXISTS trg_blog_published_at ON blog;
DROP FUNCTION IF EXISTS set_blog_published_at();
DROP TRIGGER IF EXISTS trg_blog_updated_at ON blog;
DROP TABLE IF EXISTS blog CASCADE;
