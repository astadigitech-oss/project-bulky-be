-- Records the amount confirmed by Xendit's payment.capture webhook.
ALTER TABLE xendit_payment_session
    ADD COLUMN IF NOT EXISTS captured_amount DECIMAL(15,2);
