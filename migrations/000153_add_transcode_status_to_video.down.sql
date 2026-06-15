DROP INDEX IF EXISTS idx_video_transcode_status;
ALTER TABLE video
  DROP COLUMN IF EXISTS transcode_status,
  DROP COLUMN IF EXISTS transcode_error;
