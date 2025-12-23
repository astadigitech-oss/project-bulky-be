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
