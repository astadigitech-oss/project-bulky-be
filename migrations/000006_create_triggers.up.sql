-- Trigger: Auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_kategori_produk_updated_at
    BEFORE UPDATE ON kategori_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_merek_produk_updated_at
    BEFORE UPDATE ON merek_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_kondisi_produk_updated_at
    BEFORE UPDATE ON kondisi_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_kondisi_paket_updated_at
    BEFORE UPDATE ON kondisi_paket
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sumber_produk_updated_at
    BEFORE UPDATE ON sumber_produk
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
