CREATE TABLE buyer (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    telepon VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    email_verified_at TIMESTAMP,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_buyer_username ON buyer(username);
CREATE INDEX idx_buyer_email ON buyer(email);
CREATE INDEX idx_buyer_telepon ON buyer(telepon) WHERE telepon IS NOT NULL;
CREATE INDEX idx_buyer_is_active ON buyer(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_buyer_is_verified ON buyer(is_verified);
CREATE INDEX idx_buyer_created_at ON buyer(created_at DESC);

-- Trigger for updated_at
CREATE TRIGGER update_buyer_updated_at
    BEFORE UPDATE ON buyer
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table & Column Comments
COMMENT ON TABLE buyer IS 'Menyimpan data pembeli/customer. Admin hanya memiliki akses RUD (Read, Update, Delete), buyer melakukan registrasi sendiri.';
COMMENT ON COLUMN buyer.id IS 'Primary key UUID';
COMMENT ON COLUMN buyer.nama IS 'Nama lengkap buyer';
COMMENT ON COLUMN buyer.username IS 'Username unik untuk login';
COMMENT ON COLUMN buyer.email IS 'Email unik untuk login dan notifikasi';
COMMENT ON COLUMN buyer.password IS 'Password yang sudah di-hash menggunakan bcrypt';
COMMENT ON COLUMN buyer.telepon IS 'Nomor telepon buyer (opsional)';
COMMENT ON COLUMN buyer.is_active IS 'Status aktif akun, false = akun dinonaktifkan oleh admin';
COMMENT ON COLUMN buyer.is_verified IS 'Status verifikasi email';
COMMENT ON COLUMN buyer.email_verified_at IS 'Timestamp ketika email berhasil diverifikasi';
COMMENT ON COLUMN buyer.last_login_at IS 'Timestamp login terakhir';
COMMENT ON COLUMN buyer.created_at IS 'Waktu registrasi';
COMMENT ON COLUMN buyer.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN buyer.deleted_at IS 'Soft delete timestamp';
