ALTER TABLE session_rooms
ADD CONSTRAINT session_rooms_user_id_room_type_key
UNIQUE (user_id, room_type);