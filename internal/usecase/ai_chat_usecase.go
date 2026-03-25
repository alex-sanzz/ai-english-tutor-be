package usecase

import (
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"ai-tutor-backend/internal/service"
	"context"

	"fmt"
)

type AiChatUseCase struct {
	chatService service.ChatService
	logger log.Logger
}

func NewOpenAiUseCase(chatService service.ChatService, logger log.Logger) *AiChatUseCase{
	return &AiChatUseCase{
		chatService: chatService,
		logger: logger,
	}
}

func (u *AiChatUseCase) FindRecentMessages(ctx context.Context, sessionId string, numLastMessages int32) ([]*models.Chat, error){

	return u.chatService.FindRecentMessages(ctx, sessionId, numLastMessages)

	
}


func (u *AiChatUseCase) ChatStream(ctx context.Context, sessionId string,message string, onChunk func(id string, chunk string) error, onFinish func(string) error, saveRequestMessage bool) error{
	err := u.chatService.ChatStream(ctx, sessionId, message, onChunk, onFinish, saveRequestMessage)

	if err != nil {
		return fmt.Errorf("ai chat usecase chat stream: %w", err)
	}

	return nil
	

	
}


