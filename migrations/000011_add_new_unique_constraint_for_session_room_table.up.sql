ALTER TABLE session_rooms DROP CONSTRAINT IF EXISTS session_rooms_user_id_room_type_id_unique;

ALTER TABLE session_rooms 
ADD CONSTRAINT unique_user_room_topic 
UNIQUE (user_id, room_type_id, topic);