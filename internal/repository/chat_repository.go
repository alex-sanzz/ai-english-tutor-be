package repository

import (
	"context"
	"ai-tutor-backend/internal/models"
)

type ChatRepository interface {
	Create(ctx context.Context, c *models.Chat) error
	FindByRoomType(ctx context.Context, roomType string) ([]*models.Chat, error)
	FindRecentMessages(ctx context.Context, sessionId string, numLastMessages int32) ([]*models.Chat, error)
}