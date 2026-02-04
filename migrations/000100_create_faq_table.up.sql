-- Create FAQ table
CREATE TABLE faq (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    question VARCHAR(500) NOT NULL,
    question_en VARCHAR(500) NOT NULL,
    answer TEXT NOT NULL,
    answer_en TEXT NOT NULL,
    urutan INT NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Index
CREATE INDEX idx_faq_urutan ON faq(urutan) WHERE deleted_at IS NULL;
CREATE INDEX idx_faq_is_active ON faq(is_active) WHERE deleted_at IS NULL;

-- Trigger: Auto update timestamp
CREATE TRIGGER trg_faq_updated_at
    BEFORE UPDATE ON faq
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE faq IS 'Tabel FAQ untuk accordion di frontend';
COMMENT ON COLUMN faq.question IS 'Pertanyaan dalam Bahasa Indonesia';
COMMENT ON COLUMN faq.question_en IS 'Pertanyaan dalam Bahasa Inggris';
COMMENT ON COLUMN faq.answer IS 'Jawaban dalam Bahasa Indonesia';
COMMENT ON COLUMN faq.answer_en IS 'Jawaban dalam Bahasa Inggris';
COMMENT ON COLUMN faq.urutan IS 'Urutan tampil di accordion';
