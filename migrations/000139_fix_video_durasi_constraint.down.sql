ALTER TABLE video DROP CONSTRAINT IF EXISTS chk_video_durasi_non_negative;
ALTER TABLE video ADD CONSTRAINT chk_video_durasi_positive CHECK (durasi_detik > 0);
