-- migrations/000029_create_pesanan.up.sql

-- Enum types
CREATE TYPE delivery_type AS ENUM ('PICKUP', 'DELIVEREE', 'FORWARDER');
CREATE TYPE payment_type AS ENUM ('REGULAR', 'SPLIT');
CREATE TYPE payment_status AS ENUM ('PENDING', 'PARTIAL', 'PAID', 'EXPIRED', 'FAILED', 'REFUNDED');
CREATE TYPE order_status AS ENUM ('PENDING', 'PROCESSING', 'READY', 'SHIPPED', 'COMPLETED', 'CANCELLED');

CREATE TABLE pesanan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kode VARCHAR(20) NOT NULL UNIQUE,
    buyer_id UUID NOT NULL REFERENCES buyer(id),
    delivery_type delivery_type NOT NULL,
    alamat_buyer_id UUID REFERENCES alamat_buyer(id),
    payment_type payment_type NOT NULL DEFAULT 'REGULAR',
    payment_status payment_status NOT NULL DEFAULT 'PENDING',
    order_status order_status NOT NULL DEFAULT 'PENDING',
    
    -- Biaya
    biaya_produk DECIMAL(15,2) NOT NULL DEFAULT 0,
    biaya_pengiriman DECIMAL(15,2) NOT NULL DEFAULT 0,
    biaya_ppn DECIMAL(15,2) NOT NULL DEFAULT 0,
    biaya_lainnya DECIMAL(15,2) NOT NULL DEFAULT 0,
    total DECIMAL(15,2) NOT NULL DEFAULT 0,
    
    -- Catatan
    catatan TEXT,
    catatan_admin TEXT,
    
    -- Timestamps status
    expired_at TIMESTAMP,
    paid_at TIMESTAMP,
    processed_at TIMESTAMP,
    ready_at TIMESTAMP,
    shipped_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancelled_reason TEXT,
    
    -- External IDs
    deliveree_booking_id VARCHAR(100),
    forwarder_tracking_no VARCHAR(100),
    
    -- Standard timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_alamat_required CHECK (
        (delivery_type = 'PICKUP') OR 
        (delivery_type IN ('DELIVEREE', 'FORWARDER') AND alamat_buyer_id IS NOT NULL)
    )
);

-- Indexes
CREATE INDEX idx_pesanan_kode ON pesanan(kode);
CREATE INDEX idx_pesanan_buyer_id ON pesanan(buyer_id);
CREATE INDEX idx_pesanan_order_status ON pesanan(order_status);
CREATE INDEX idx_pesanan_payment_status ON pesanan(payment_status);
CREATE INDEX idx_pesanan_delivery_type ON pesanan(delivery_type);
CREATE INDEX idx_pesanan_created_at ON pesanan(created_at DESC);
CREATE INDEX idx_pesanan_expired_at ON pesanan(expired_at) WHERE payment_status = 'PENDING';

-- Trigger updated_at
CREATE TRIGGER trg_pesanan_updated_at
    BEFORE UPDATE ON pesanan
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Function to generate order code
CREATE OR REPLACE FUNCTION generate_order_code()
RETURNS TRIGGER AS $$
DECLARE
    v_date TEXT;
    v_sequence INT;
    v_code TEXT;
BEGIN
    v_date := TO_CHAR(NOW(), 'YYYYMMDD');
    
    SELECT COALESCE(MAX(
        CAST(SUBSTRING(kode FROM 14 FOR 4) AS INT)
    ), 0) + 1 INTO v_sequence
    FROM pesanan
    WHERE kode LIKE 'ORD-' || v_date || '-%';
    
    v_code := 'ORD-' || v_date || '-' || LPAD(v_sequence::TEXT, 4, '0');
    NEW.kode := v_code;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_generate_order_code
    BEFORE INSERT ON pesanan
    FOR EACH ROW
    WHEN (NEW.kode IS NULL)
    EXECUTE FUNCTION generate_order_code();

COMMENT ON TABLE pesanan IS 'Tabel pesanan. Support pickup, deliveree, forwarder. Support regular & split payment.';
