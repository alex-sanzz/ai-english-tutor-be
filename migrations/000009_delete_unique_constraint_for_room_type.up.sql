ALTER TABLE session_rooms DROP CONSTRAINT IF EXISTS session_rooms_user_id_room_type_key;

ALTER TABLE session_rooms ADD CONSTRAINT session_rooms_user_id_room_type_id_unique UNIQUE(user_id, room_type_id);