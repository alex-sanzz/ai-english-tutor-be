CREATE TABLE room_types (type VARCHAR(255) PRIMARY KEY);

ALTER TABLE session_rooms DROP COLUMN room_type;

ALTER TABLE session_rooms ADD COLUMN room_type_id VARCHAR(255), ADD COLUMN topic TEXT;

ALTER TABLE session_rooms ADD CONSTRAINT fk_session_room_room_type FOREIGN KEY (room_type_id) REFERENCES room_types(type);
