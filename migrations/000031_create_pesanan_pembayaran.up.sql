-- migrations/000031_create_pesanan_pembayaran.up.sql

CREATE TABLE pesanan_pembayaran (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pesanan_id UUID NOT NULL REFERENCES pesanan(id) ON DELETE CASCADE,
    buyer_id UUID NOT NULL REFERENCES buyer(id),
    metode_pembayaran_id UUID REFERENCES metode_pembayaran(id),
    jumlah DECIMAL(15,2) NOT NULL,
    status payment_status NOT NULL DEFAULT 'PENDING',
    
    -- Xendit fields
    xendit_invoice_id VARCHAR(100),
    xendit_external_id VARCHAR(100) UNIQUE,
    xendit_payment_url TEXT,
    xendit_payment_method VARCHAR(50),
    
    -- Timestamps
    expired_at TIMESTAMP,
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_pesanan_pembayaran_pesanan_id ON pesanan_pembayaran(pesanan_id);
CREATE INDEX idx_pesanan_pembayaran_buyer_id ON pesanan_pembayaran(buyer_id);
CREATE INDEX idx_pesanan_pembayaran_status ON pesanan_pembayaran(status);
CREATE INDEX idx_pesanan_pembayaran_xendit_invoice_id ON pesanan_pembayaran(xendit_invoice_id);
CREATE INDEX idx_pesanan_pembayaran_xendit_external_id ON pesanan_pembayaran(xendit_external_id);

CREATE TRIGGER trg_pesanan_pembayaran_updated_at
    BEFORE UPDATE ON pesanan_pembayaran
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Constraint: Max 5 payments per order (split payment limit)
CREATE OR REPLACE FUNCTION check_max_split_payment()
RETURNS TRIGGER AS $$
DECLARE
    v_count INT;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM pesanan_pembayaran
    WHERE pesanan_id = NEW.pesanan_id;
    
    IF v_count >= 5 THEN
        RAISE EXCEPTION 'Maksimal 5 pembayaran per pesanan (split payment limit)';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_max_split_payment
    BEFORE INSERT ON pesanan_pembayaran
    FOR EACH ROW
    EXECUTE FUNCTION check_max_split_payment();

COMMENT ON TABLE pesanan_pembayaran IS 'Pembayaran per pesanan. Support split payment max 5 orang.';
