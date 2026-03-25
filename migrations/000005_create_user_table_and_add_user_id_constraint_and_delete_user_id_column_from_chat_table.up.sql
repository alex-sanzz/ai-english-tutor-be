CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    google_sub TEXT UNIQUE NOT NULL ,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE chats DROP COLUMN user_id;

ALTER TABLE session_rooms DROP COLUMN user_id;

ALTER TABLE session_rooms ADD COLUMN user_id UUID;

ALTER TABLE session_rooms ADD CONSTRAINT fk_session_room_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE;



