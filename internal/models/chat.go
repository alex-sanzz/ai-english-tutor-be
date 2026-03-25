package models

import "time"

type Chat struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	SessionRoomId string `json:"session_room_id"`
	Role string `json:"role"`
	ReceivedTime    time.Time `json:"time"`
	CreatedAt time.Time
	UpdatedAt time.Time
}