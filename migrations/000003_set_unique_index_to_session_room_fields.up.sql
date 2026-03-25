
/* UNIQUE INDEX means index those specific columns, then make it unique */
/* Create an index for speed AND use it to guarantee that values are unique */
CREATE UNIQUE INDEX idx_unique_user_id_and_room_type ON session_rooms(user_id, room_type);