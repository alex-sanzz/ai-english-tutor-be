package dto

type CreateSessionRoomRequest struct {
	RoomType string `json:"room_type"`
	DetailedContext string `json:"detailed_context"`
	Topic string  `json:"topic"`
	Icon string `json:"icon"`
}