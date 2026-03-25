ALTER TABLE session_rooms ADD CONSTRAINT session_rooms_user_id_room_type_id_unique UNIQUE(user_id, room_type_id);

ALTER TABLE session_rooms DROP CONSTRAINT IF EXISTS unique_user_room_topic;
