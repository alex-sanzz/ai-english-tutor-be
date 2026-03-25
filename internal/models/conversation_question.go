package models

import (
	"time"

)

type ConversationQuestion struct {
	ID           string `json:"id"`
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	ReviewResult string `json:"review_result"`
	SessionRoomId string `json:"session_room_id"`
	AnsweredAt time.Time `json:"answered_at"`
	CreatedAt    time.Time `json:"created_at"`
}