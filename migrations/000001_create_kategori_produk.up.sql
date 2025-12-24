-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabel: kategori_produk
CREATE TABLE kategori_produk (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE,
    deskripsi TEXT,
    icon_url VARCHAR(500),
    memiliki_kondisi_tambahan BOOLEAN DEFAULT false,
    tipe_kondisi_tambahan VARCHAR(20) CHECK (tipe_kondisi_tambahan IN ('gambar', 'teks')),
    gambar_kondisi_url VARCHAR(500),
    teks_kondisi TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_kategori_produk_slug ON kategori_produk(slug);
CREATE INDEX idx_kategori_produk_is_active ON kategori_produk(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_kategori_produk_deleted_at ON kategori_produk(deleted_at);

-- Table & Column Comments
COMMENT ON TABLE kategori_produk IS 'Menyimpan data kategori produk dengan fitur kondisi tambahan (gambar/teks)';
COMMENT ON COLUMN kategori_produk.id IS 'Primary key UUID';
COMMENT ON COLUMN kategori_produk.nama IS 'Nama kategori';
COMMENT ON COLUMN kategori_produk.slug IS 'URL-friendly identifier';
COMMENT ON COLUMN kategori_produk.deskripsi IS 'Deskripsi kategori';
COMMENT ON COLUMN kategori_produk.icon_url IS 'URL icon kategori';
COMMENT ON COLUMN kategori_produk.memiliki_kondisi_tambahan IS 'Flag apakah kategori memiliki kondisi tambahan';
COMMENT ON COLUMN kategori_produk.tipe_kondisi_tambahan IS 'Tipe kondisi tambahan (gambar/teks)';
COMMENT ON COLUMN kategori_produk.gambar_kondisi_url IS 'URL gambar kondisi tambahan';
COMMENT ON COLUMN kategori_produk.teks_kondisi IS 'Teks kondisi tambahan';
COMMENT ON COLUMN kategori_produk.is_active IS 'Status aktif';
COMMENT ON COLUMN kategori_produk.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN kategori_produk.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN kategori_produk.deleted_at IS 'Soft delete timestamp';
