ALTER TABLE pesanan ADD COLUMN booking_error TEXT;

COMMENT ON COLUMN pesanan.booking_error IS
    'Error message dari Deliveree/Forwarder API jika booking gagal. NULL jika sukses atau belum di-booking.';
