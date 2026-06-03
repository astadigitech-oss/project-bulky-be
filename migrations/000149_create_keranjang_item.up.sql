-- migrations/000149_create_keranjang_item.up.sql

CREATE TABLE keranjang_item (
    id              UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    keranjang_id    UUID        NOT NULL REFERENCES keranjang(id) ON DELETE CASCADE,
    produk_id       UUID        NOT NULL REFERENCES produk(id) ON DELETE CASCADE,
    quantity        INTEGER     NOT NULL DEFAULT 1 CHECK (quantity > 0),
    is_selected     BOOLEAN     NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT keranjang_item_keranjang_produk_unique UNIQUE (keranjang_id, produk_id)
);

CREATE INDEX idx_keranjang_item_keranjang_id ON keranjang_item(keranjang_id);
CREATE INDEX idx_keranjang_item_produk_id ON keranjang_item(produk_id);

CREATE TRIGGER update_keranjang_item_updated_at
    BEFORE UPDATE ON keranjang_item
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE keranjang_item IS 'Item dalam keranjang belanja. Satu baris = satu produk. Tidak ada duplikat produk dalam satu cart (UNIQUE keranjang_id + produk_id).';
COMMENT ON COLUMN keranjang_item.quantity IS 'Jumlah item. Untuk Paletbox selalu 1 karena tiap produk adalah unit unik.';
COMMENT ON COLUMN keranjang_item.is_selected IS 'Status centang item untuk checkout. Default true saat item ditambahkan. False = item di-uncheck oleh buyer.';
