CREATE TABLE admin (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_admin_email ON admin(email);
CREATE INDEX idx_admin_is_active ON admin(is_active) WHERE deleted_at IS NULL;

-- Trigger for updated_at
CREATE TRIGGER update_admin_updated_at
    BEFORE UPDATE ON admin
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table & Column Comments
COMMENT ON TABLE admin IS 'Menyimpan data administrator sistem';
COMMENT ON COLUMN admin.id IS 'Primary key UUID';
COMMENT ON COLUMN admin.nama IS 'Nama lengkap admin';
COMMENT ON COLUMN admin.email IS 'Email untuk login (unique)';
COMMENT ON COLUMN admin.password IS 'Password (hashed bcrypt)';
COMMENT ON COLUMN admin.is_active IS 'Status aktif akun';
COMMENT ON COLUMN admin.last_login_at IS 'Waktu login terakhir';
COMMENT ON COLUMN admin.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN admin.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN admin.deleted_at IS 'Soft delete timestamp';
