-- migrations/000032_create_pesanan_status_history.up.sql

CREATE TYPE status_history_type AS ENUM ('ORDER', 'PAYMENT');

CREATE TABLE pesanan_status_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pesanan_id UUID NOT NULL REFERENCES pesanan(id) ON DELETE CASCADE,
    status_from VARCHAR(20),
    status_to VARCHAR(20) NOT NULL,
    status_type status_history_type NOT NULL,
    changed_by UUID REFERENCES admin(id),
    note TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_pesanan_status_history_pesanan_id ON pesanan_status_history(pesanan_id);
CREATE INDEX idx_pesanan_status_history_created_at ON pesanan_status_history(created_at);

-- Trigger to auto-log status changes
CREATE OR REPLACE FUNCTION log_pesanan_status_change()
RETURNS TRIGGER AS $$
BEGIN
    -- Log order status change
    IF OLD.order_status IS DISTINCT FROM NEW.order_status THEN
        INSERT INTO pesanan_status_history (pesanan_id, status_from, status_to, status_type)
        VALUES (NEW.id, OLD.order_status::TEXT, NEW.order_status::TEXT, 'ORDER');
    END IF;
    
    -- Log payment status change
    IF OLD.payment_status IS DISTINCT FROM NEW.payment_status THEN
        INSERT INTO pesanan_status_history (pesanan_id, status_from, status_to, status_type)
        VALUES (NEW.id, OLD.payment_status::TEXT, NEW.payment_status::TEXT, 'PAYMENT');
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_log_pesanan_status_change
    AFTER UPDATE ON pesanan
    FOR EACH ROW
    EXECUTE FUNCTION log_pesanan_status_change();

COMMENT ON TABLE pesanan_status_history IS 'Log perubahan status pesanan untuk audit trail.';
