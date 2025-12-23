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
