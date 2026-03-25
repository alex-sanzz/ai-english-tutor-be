CREATE TABLE conversation_questions (
    id UUID PRIMARY KEY default gen_random_uuid(),
    question TEXT NOT NULL,
    answer TEXT,
    review_result TEXT,
    session_room_id UUID,
    answered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE conversation_questions ADD CONSTRAINT fk_conversation_question_session_room FOREIGN KEY(session_room_id) REFERENCES session_rooms(id) ON DELETE CASCADE;