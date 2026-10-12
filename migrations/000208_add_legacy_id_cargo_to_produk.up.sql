ALTER TABLE produk
    ADD COLUMN legacy_id_cargo BIGINT;

CREATE UNIQUE INDEX idx_produk_legacy_id_cargo
    ON produk(legacy_id_cargo);
