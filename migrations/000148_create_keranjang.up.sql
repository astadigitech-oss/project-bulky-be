-- migrations/000148_create_keranjang.up.sql

CREATE TABLE keranjang (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    buyer_id    UUID        NOT NULL REFERENCES buyer(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT keranjang_buyer_id_unique UNIQUE (buyer_id)
);

CREATE INDEX idx_keranjang_buyer_id ON keranjang(buyer_id);

CREATE TRIGGER update_keranjang_updated_at
    BEFORE UPDATE ON keranjang
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE keranjang IS 'Keranjang belanja buyer. Satu buyer satu cart aktif (UNIQUE buyer_id). Cart dibuat otomatis saat buyer pertama menambahkan produk.';
COMMENT ON COLUMN keranjang.buyer_id IS 'FK ke tabel buyer. Constraint UNIQUE memastikan hanya ada 1 cart per buyer.';
