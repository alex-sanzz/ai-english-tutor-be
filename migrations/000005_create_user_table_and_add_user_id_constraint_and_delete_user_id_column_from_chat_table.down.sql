DELETE TABLE users;

ALTER TABLE chats ADD COLUMN user_id TEXT;

ALTER TABLE session_rooms DROP COLUMN user_id;

ALTER TABLE session_rooms ADD COLUMN user_id TEXT;

ALTER TABLE session_rooms DROP CONSTRAINT fk_session_room_user;
