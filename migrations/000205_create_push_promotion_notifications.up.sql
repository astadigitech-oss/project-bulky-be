CREATE TABLE push_promotion_notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama VARCHAR(120) NOT NULL,
    title_id VARCHAR(160) NOT NULL,
    body_id TEXT NOT NULL,
    title_en VARCHAR(160) NOT NULL,
    body_en TEXT NOT NULL,
    deep_link VARCHAR(500) NOT NULL DEFAULT '/products',
    scheduled_at TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT'
        CHECK (status IN ('DRAFT', 'SCHEDULED', 'SENDING', 'SENT', 'FAILED')),
    sent_at TIMESTAMPTZ,
    success_count INTEGER NOT NULL DEFAULT 0,
    failure_count INTEGER NOT NULL DEFAULT 0,
    last_error TEXT,
    created_by UUID REFERENCES admin(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES admin(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_push_promotion_scheduled_at
        CHECK (status <> 'SCHEDULED' OR scheduled_at IS NOT NULL)
);

CREATE INDEX idx_push_promotion_notifications_due
    ON push_promotion_notifications (scheduled_at, id)
    WHERE status = 'SCHEDULED';

CREATE INDEX idx_push_promotion_notifications_created
    ON push_promotion_notifications (created_at DESC);
