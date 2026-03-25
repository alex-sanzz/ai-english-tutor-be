package infrastructure

import (
	"ai-tutor-backend/internal/models"
	"context"
)

type AiChatClient interface {
	ChatStream(ctx context.Context, messages []*models.Chat, onChunk func(string) error, onFinish func(string) error) error
	AskQuestion(ctx context.Context, question string) (string, error)
}
