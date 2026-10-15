CREATE TABLE warehouse_visit_dates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    warehouse_id UUID NOT NULL REFERENCES warehouse(id) ON DELETE RESTRICT,
    visit_date DATE NOT NULL,
    is_open BOOLEAN NOT NULL DEFAULT TRUE,
    crew_count INTEGER NOT NULL DEFAULT 0 CHECK (crew_count >= 0),
    visit_hours_note TEXT,
    created_by UUID REFERENCES admin(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES admin(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_warehouse_visit_dates_warehouse_date UNIQUE (warehouse_id, visit_date),
    CONSTRAINT uq_warehouse_visit_dates_id_date UNIQUE (id, visit_date),
    CONSTRAINT chk_warehouse_visit_dates_open_has_crew CHECK (NOT is_open OR crew_count > 0)
);

CREATE INDEX idx_warehouse_visit_dates_date ON warehouse_visit_dates(visit_date, is_open);

CREATE TRIGGER trg_warehouse_visit_dates_updated_at
    BEFORE UPDATE ON warehouse_visit_dates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE warehouse_visit_bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    visit_date_id UUID NOT NULL,
    visit_date DATE NOT NULL,
    buyer_id UUID NOT NULL REFERENCES buyer(id) ON DELETE RESTRICT,
    status VARCHAR(20) NOT NULL DEFAULT 'CONFIRMED'
        CHECK (status IN ('CONFIRMED', 'CANCELLED', 'ATTENDED', 'NO_SHOW')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_warehouse_visit_bookings_date
        FOREIGN KEY (visit_date_id, visit_date)
        REFERENCES warehouse_visit_dates(id, visit_date) ON DELETE RESTRICT
);

CREATE UNIQUE INDEX uq_warehouse_visit_bookings_buyer_date_active
    ON warehouse_visit_bookings(buyer_id, visit_date)
    WHERE status <> 'CANCELLED';

CREATE INDEX idx_warehouse_visit_bookings_date_status
    ON warehouse_visit_bookings(visit_date_id, status);

CREATE TRIGGER trg_warehouse_visit_bookings_updated_at
    BEFORE UPDATE ON warehouse_visit_bookings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE warehouse_visit_dates IS 'Tanggal kunjungan gudang yang dibuka Admin; kapasitas = crew_count x 2 buyer.';
COMMENT ON COLUMN warehouse_visit_dates.visit_hours_note IS 'Catatan jam kunjungan untuk informasi; bukan slot reservasi.';
COMMENT ON TABLE warehouse_visit_bookings IS 'Reservasi kunjungan gudang per buyer dan tanggal.';
