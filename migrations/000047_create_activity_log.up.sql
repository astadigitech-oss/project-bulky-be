-- migrations/000047_create_activity_log.up.sql

CREATE TABLE activity_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_type VARCHAR(20) NOT NULL CHECK (user_type IN ('ADMIN', 'BUYER', 'SYSTEM')),
    user_id UUID,
    action VARCHAR(50) NOT NULL,
    modul VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50),
    entity_id UUID,
    deskripsi TEXT NOT NULL,
    old_data JSONB,
    new_data JSONB,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_activity_log_user ON activity_log(user_type, user_id);
CREATE INDEX idx_activity_log_action ON activity_log(action);
CREATE INDEX idx_activity_log_modul ON activity_log(modul);
CREATE INDEX idx_activity_log_entity ON activity_log(entity_type, entity_id);
CREATE INDEX idx_activity_log_created_at ON activity_log(created_at DESC);

COMMENT ON TABLE activity_log IS 'Audit trail untuk setiap aktivitas penting di sistem';
COMMENT ON COLUMN activity_log.old_data IS 'Data sebelum perubahan (untuk UPDATE/DELETE)';
COMMENT ON COLUMN activity_log.new_data IS 'Data setelah perubahan (untuk CREATE/UPDATE)';
