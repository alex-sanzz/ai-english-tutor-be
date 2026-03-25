CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE session_rooms(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message TEXT NOT NULL,
    user_id TEXT NOT NULL,
    session_room_id UUID NOT NULL,
    received_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    /* -- This basically tells when a row from session rooms is deleted, then just delete the corresponding chat rows */
    CONSTRAINT fk_session_room FOREIGN KEY(session_room_id) REFERENCES session_rooms(id) ON DELETE CASCADE
);

