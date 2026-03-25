package dto

type CreateSessionRoomRequest struct {
	RoomType string `json:"room_type"`
	Topic string  `json:"topic"`
	Icon string `json:"icon"`
}