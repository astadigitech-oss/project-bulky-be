CREATE TABLE forwarder_city_mapping (
    id                  SERIAL PRIMARY KEY,
    kota_pattern        VARCHAR(100) NOT NULL,
    forwarder_city_id   INT NOT NULL,
    forwarder_city_name VARCHAR(100) NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_forwarder_city_mapping_pattern
    ON forwarder_city_mapping (kota_pattern);

CREATE INDEX idx_forwarder_city_mapping_city_id
    ON forwarder_city_mapping (forwarder_city_id);

COMMENT ON TABLE forwarder_city_mapping IS
    'Mapping kota teks (Google Maps) ke Forwarder.ai city ID. Di-seed dari API citylist Forwarder.';
COMMENT ON COLUMN forwarder_city_mapping.kota_pattern IS
    'Kota yang sudah dinormalisasi: lowercase, stripped prefix "kota"/"kabupaten". Contoh: "yogyakarta", "bogor".';
