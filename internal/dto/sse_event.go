package dto

type SseEvent struct {
	MessageId 	  string `json:"message_id"`
	Message       string `json:"message"`
	Role        string `json:"role"`
}
