CREATE TABLE kondisi_paket (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(100) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE,
    deskripsi TEXT,
    urutan INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_kondisi_paket_slug ON kondisi_paket(slug);
CREATE INDEX idx_kondisi_paket_urutan ON kondisi_paket(urutan);

-- Table & Column Comments
COMMENT ON TABLE kondisi_paket IS 'Menyimpan kondisi kelengkapan paket (Fullset, Tanpa Dus, Unit Only, dll)';
COMMENT ON COLUMN kondisi_paket.id IS 'Primary key UUID';
COMMENT ON COLUMN kondisi_paket.nama IS 'Nama kondisi paket';
COMMENT ON COLUMN kondisi_paket.slug IS 'URL-friendly identifier';
COMMENT ON COLUMN kondisi_paket.deskripsi IS 'Deskripsi kondisi paket';
COMMENT ON COLUMN kondisi_paket.urutan IS 'Urutan tampilan';
COMMENT ON COLUMN kondisi_paket.is_active IS 'Status aktif';
COMMENT ON COLUMN kondisi_paket.created_at IS 'Waktu dibuat';
COMMENT ON COLUMN kondisi_paket.updated_at IS 'Waktu terakhir diupdate';
COMMENT ON COLUMN kondisi_paket.deleted_at IS 'Soft delete timestamp';
