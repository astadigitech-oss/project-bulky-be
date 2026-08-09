-- Hapus kolom lock anti double-booking.
ALTER TABLE pesanan
    DROP COLUMN IF EXISTS booking_lock_at;
