-- =====================================================
-- TABEL: formulir_partai_besar_submission
-- =====================================================
-- History submission formulir partai besar
-- =====================================================

CREATE TABLE formulir_partai_besar_submission (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    buyer_id UUID REFERENCES buyer(id),
    nama VARCHAR(100) NOT NULL,
    telepon VARCHAR(20) NOT NULL,
    alamat TEXT NOT NULL,
    anggaran_id UUID REFERENCES formulir_partai_besar_anggaran(id),
    kategori_ids TEXT NOT NULL DEFAULT '[]',
    email_sent BOOLEAN DEFAULT false,
    email_sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_formulir_submission_buyer_id ON formulir_partai_besar_submission(buyer_id);
CREATE INDEX idx_formulir_submission_created_at ON formulir_partai_besar_submission(created_at DESC);
CREATE INDEX idx_formulir_submission_email_sent ON formulir_partai_besar_submission(email_sent);

COMMENT ON TABLE formulir_partai_besar_submission IS 'History submission formulir partai besar';
COMMENT ON COLUMN formulir_partai_besar_submission.kategori_ids IS 'JSON array UUID kategori yang dipilih';
