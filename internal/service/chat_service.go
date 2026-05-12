package service

import (
	"ai-tutor-backend/internal/models"
	"context"
)

type ChatService interface {
	ChatStream(ctx context.Context, sessionId string, message string, systemPrompt string, onChunk func(id string, chunk string) error, onFinish func(string) error, saveRequestMessage bool) error
	FindRecentMessages(ctx context.Context, sessionId string, numLastMessages int32) ([]*models.Chat, error)
}