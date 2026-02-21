ALTER TABLE video DROP CONSTRAINT IF EXISTS chk_video_durasi_positive;
ALTER TABLE video ADD CONSTRAINT chk_video_durasi_non_negative CHECK (durasi_detik >= 0);
