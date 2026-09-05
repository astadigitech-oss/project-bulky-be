CREATE TABLE xendit_payment_session (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pesanan_pembayaran_id UUID NOT NULL
        REFERENCES pesanan_pembayaran(id) ON DELETE CASCADE,
    reference_id VARCHAR(100) NOT NULL,
    payment_session_id VARCHAR(100) NOT NULL UNIQUE,
    payment_request_id VARCHAR(100),
    payment_id VARCHAR(100) UNIQUE,
    payment_link_url TEXT NOT NULL,
    channel_code VARCHAR(50) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'IDR',
    status VARCHAR(30) NOT NULL CHECK (status IN ('ACTIVE', 'COMPLETED', 'EXPIRED', 'CANCELED', 'FAILED')),
    failure_code VARCHAR(100),
    failure_message TEXT,
    expires_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uq_xendit_payment_session_active_per_payment
    ON xendit_payment_session (pesanan_pembayaran_id)
    WHERE status = 'ACTIVE';

CREATE INDEX idx_xendit_payment_session_reference_id
    ON xendit_payment_session (reference_id);

CREATE INDEX idx_xendit_payment_session_payment_status
    ON xendit_payment_session (pesanan_pembayaran_id, created_at DESC);

CREATE TABLE xendit_webhook_event (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event VARCHAR(100) NOT NULL,
    payment_session_id VARCHAR(100),
    payment_id VARCHAR(100),
    reference_id VARCHAR(100),
    payload JSONB NOT NULL,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_xendit_webhook_event_session_id
    ON xendit_webhook_event (payment_session_id)
    WHERE payment_session_id IS NOT NULL;

CREATE INDEX idx_xendit_webhook_event_payment_id
    ON xendit_webhook_event (payment_id)
    WHERE payment_id IS NOT NULL;

CREATE INDEX idx_xendit_webhook_event_reference_id
    ON xendit_webhook_event (reference_id)
    WHERE reference_id IS NOT NULL;

CREATE UNIQUE INDEX uq_xendit_webhook_event_resource
    ON xendit_webhook_event (
        event,
        COALESCE(payment_session_id, ''),
        COALESCE(payment_id, ''),
        COALESCE(reference_id, '')
    );
