CREATE TABLE subcriptions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    purchase_token TEXT NOT NULL,
    purchase_time TIMESTAMPTZ NOT NULL,
    expiry_at TIMESTAMPTZ NOT NULL,
    auto_renewing BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE subcriptions ADD CONSTRAINT fk_subcription_user_table FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE;