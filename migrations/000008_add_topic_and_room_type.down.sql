ALTER TABLE session_rooms ADD COLUMN room_type VARCHAR(255);

ALTER TABLE session_rooms DROP FOREIGN KEY fk_session_rooms_room_type, DROP COLUMN room_type_id, DROP COLUMN topic;

DROP TABLE room_types;