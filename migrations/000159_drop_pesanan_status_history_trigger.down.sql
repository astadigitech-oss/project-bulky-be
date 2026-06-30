-- Recreate trigger and function for auto-logging status changes
CREATE OR REPLACE FUNCTION log_pesanan_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.order_status IS DISTINCT FROM NEW.order_status THEN
        INSERT INTO pesanan_status_history (pesanan_id, status_from, status_to, status_type)
        VALUES (NEW.id, OLD.order_status::TEXT, NEW.order_status::TEXT, 'ORDER');
    END IF;

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
