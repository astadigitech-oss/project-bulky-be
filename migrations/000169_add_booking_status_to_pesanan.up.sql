-- Tambah kolom status booking & tracking URL dari provider pengiriman
-- (Deliveree/Forwarder) di tabel pesanan. Paritas dengan BE lama (panel-bulky)
-- yang menyimpan booking_status dan tracking_url di tabel order_shipping.
--
-- booking_status diisi oleh webhook provider (mis. delivery_in_progress,
-- delivery_completed) atau status booking terakhir saat admin cek tracking.
-- tracking_url adalah link live-tracking yang bisa dibagikan ke buyer.

ALTER TABLE pesanan
    ADD COLUMN IF NOT EXISTS booking_status VARCHAR(50),
    ADD COLUMN IF NOT EXISTS tracking_url TEXT;

COMMENT ON COLUMN pesanan.booking_status IS 'Status terakhir dari provider pengiriman (Deliveree/Forwarder), diperbarui via webhook';
COMMENT ON COLUMN pesanan.tracking_url IS 'Link live tracking pengiriman dari provider';
