ALTER TABLE video
  ADD COLUMN transcode_status VARCHAR(20)
    NOT NULL
    DEFAULT 'pending'
    CHECK (transcode_status IN ('pending', 'processing', 'ready', 'failed')),
  ADD COLUMN transcode_error TEXT;

CREATE INDEX idx_video_transcode_status ON video(transcode_status)
  WHERE deleted_at IS NULL;

COMMENT ON COLUMN video.transcode_status IS 'Status proses transcode: pending | processing | ready | failed';
COMMENT ON COLUMN video.transcode_error IS 'Pesan error dari FFmpeg jika transcode gagal';
