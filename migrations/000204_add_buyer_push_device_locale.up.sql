ALTER TABLE buyer_push_devices
    ADD COLUMN locale VARCHAR(5) NOT NULL DEFAULT 'id'
    CHECK (locale IN ('id', 'en'));

COMMENT ON COLUMN buyer_push_devices.locale IS
    'Preferred locale used to localize buyer push notification content.';
