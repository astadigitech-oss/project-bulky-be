-- migrations/000030_create_pesanan_item.up.sql

CREATE TABLE pesanan_item (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pesanan_id UUID NOT NULL REFERENCES pesanan(id) ON DELETE CASCADE,
    produk_id UUID NOT NULL REFERENCES produk(id),
    nama_produk VARCHAR(200) NOT NULL,
    sku VARCHAR(50),
    qty INT NOT NULL CHECK (qty > 0),
    harga_satuan DECIMAL(15,2) NOT NULL,
    diskon_satuan DECIMAL(15,2) DEFAULT 0,
    subtotal DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_pesanan_item_pesanan_id ON pesanan_item(pesanan_id);
CREATE INDEX idx_pesanan_item_produk_id ON pesanan_item(produk_id);

CREATE TRIGGER trg_pesanan_item_updated_at
    BEFORE UPDATE ON pesanan_item
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger to calculate subtotal
CREATE OR REPLACE FUNCTION calculate_pesanan_item_subtotal()
RETURNS TRIGGER AS $$
BEGIN
    NEW.subtotal := NEW.qty * (NEW.harga_satuan - COALESCE(NEW.diskon_satuan, 0));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_calculate_subtotal
    BEFORE INSERT OR UPDATE ON pesanan_item
    FOR EACH ROW
    EXECUTE FUNCTION calculate_pesanan_item_subtotal();

COMMENT ON TABLE pesanan_item IS 'Detail item dalam pesanan. Harga di-snapshot saat checkout.';
