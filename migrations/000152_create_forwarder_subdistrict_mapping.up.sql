CREATE TABLE forwarder_subdistrict_mapping (
    id                          SERIAL PRIMARY KEY,
    kecamatan_pattern           VARCHAR(100) NOT NULL,
    forwarder_city_id           INT NOT NULL,
    forwarder_subdistrict_id    INT NOT NULL,
    forwarder_subdistrict_name  VARCHAR(100) NOT NULL,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_forwarder_subdistrict_mapping_pattern_city
    ON forwarder_subdistrict_mapping (kecamatan_pattern, forwarder_city_id);

CREATE INDEX idx_forwarder_subdistrict_mapping_subdistrict_id
    ON forwarder_subdistrict_mapping (forwarder_subdistrict_id);

COMMENT ON TABLE forwarder_subdistrict_mapping IS
    'Mapping kecamatan teks (Google Maps) ke Forwarder.ai subdistrict ID. Di-seed dari API subdistrictlist Forwarder.';
COMMENT ON COLUMN forwarder_subdistrict_mapping.kecamatan_pattern IS
    'Kecamatan yang sudah dinormalisasi: lowercase. Contoh: "kotagede", "cibinong".';
