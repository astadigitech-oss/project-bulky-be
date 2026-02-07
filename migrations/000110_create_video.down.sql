DROP TRIGGER IF EXISTS trg_video_published_at ON video;
DROP FUNCTION IF EXISTS set_video_published_at();
DROP TRIGGER IF EXISTS trg_video_updated_at ON video;
DROP TABLE IF EXISTS video CASCADE;
