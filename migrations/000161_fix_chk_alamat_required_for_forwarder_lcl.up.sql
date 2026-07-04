ALTER TABLE pesanan DROP CONSTRAINT chk_alamat_required;

ALTER TABLE pesanan ADD CONSTRAINT chk_alamat_required CHECK (
    (delivery_type = 'PICKUP') OR
    (delivery_type IN ('DELIVEREE', 'FORWARDER', 'FORWARDER_LCL') AND alamat_buyer_id IS NOT NULL)
);
