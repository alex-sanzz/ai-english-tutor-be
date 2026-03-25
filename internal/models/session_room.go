package models

import (
	"time"

	// "github.com/google/uuid"
)

type SessionRoom struct {
	ID        string `json:"id"`
	RoomType  string `json:"room_type"`
	UserId  string `json:"user_id"`
	Topic string `json:"topic"`
	Icon string `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
}