-- Rollback: Re-enable triggers

-- kategori_produk
CREATE TRIGGER trigger_rewrite_slug_on_delete_kategori_produk
    BEFORE UPDATE ON kategori_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- merek_produk
CREATE TRIGGER trigger_rewrite_slug_on_delete_merek_produk
    BEFORE UPDATE ON merek_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- kondisi_produk
CREATE TRIGGER trigger_rewrite_slug_on_delete_kondisi_produk
    BEFORE UPDATE ON kondisi_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- kondisi_paket
CREATE TRIGGER trigger_rewrite_slug_on_delete_kondisi_paket
    BEFORE UPDATE ON kondisi_paket
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- sumber_produk
CREATE TRIGGER trigger_rewrite_slug_on_delete_sumber_produk
    BEFORE UPDATE ON sumber_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- warehouse
CREATE TRIGGER trigger_rewrite_slug_on_delete_warehouse
    BEFORE UPDATE ON warehouse
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- tipe_produk
CREATE TRIGGER trigger_rewrite_slug_on_delete_tipe_produk
    BEFORE UPDATE ON tipe_produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- produk
CREATE TRIGGER trigger_rewrite_slug_on_delete_produk
    BEFORE UPDATE ON produk
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- dokumen_kebijakan
CREATE TRIGGER trigger_rewrite_slug_on_delete_dokumen_kebijakan
    BEFORE UPDATE ON dokumen_kebijakan
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();

-- disclaimer
CREATE TRIGGER trigger_rewrite_slug_on_delete_disclaimer
    BEFORE UPDATE ON disclaimer
    FOR EACH ROW
    EXECUTE FUNCTION rewrite_slug_on_delete();
