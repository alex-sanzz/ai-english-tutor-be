package service

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/infrastructure"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"ai-tutor-backend/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

)

type chatService struct {
	aiChatClient infrastructure.AiChatClient
	chatRepo repository.ChatRepository
	logger log.Logger
}

func NewChatService(client infrastructure.AiChatClient, chatRepo repository.ChatRepository, logger log.Logger) *chatService {
	return &chatService{
		aiChatClient: client,
		chatRepo: chatRepo,
		logger: logger,
	}
}

func (s *chatService) ChatStream(ctx context.Context, sessionId string, message string, onChunk func(id string, chunk string) error, onFinish func(string) error, saveRequestMessage bool) error {
	gotMessageTime := time.Now()
	messageId := uuid.New().String()
	chats, err := s.chatRepo.FindRecentMessages(ctx, sessionId, 10)
	if err != nil {
		return fmt.Errorf("chat service chat stream error: %w", err)
	}
	messages := make([]*models.Chat, 0, len(chats) + 1)

	for _, chat := range chats {
		messages = append(messages, chat)
	}

	messages = append(messages, &models.Chat{
		Message: message,
		SessionRoomId: sessionId,
		Role: "user",
		ReceivedTime: gotMessageTime,
		CreatedAt: gotMessageTime,
		UpdatedAt: gotMessageTime,
	})

	// Implementation of streaming chat logic goes here
	return s.aiChatClient.ChatStream(ctx, messages, func(chunk string) error {
		// Handle each chunk of the streamed response
		onChunk(messageId, chunk)
		return nil
	}, func(fullSentences string) error {

		if(saveRequestMessage){
			// Handle completion of the stream
			err := s.chatRepo.Create(ctx, &models.Chat{
				Message: message,
				SessionRoomId: sessionId,
				Role: "user",
				ReceivedTime: gotMessageTime,
				CreatedAt: gotMessageTime,
				UpdatedAt: gotMessageTime,
			})

			if err != nil {	
				return apperr.Internal(fmt.Errorf("chat service chat stream on finish: %w", err))
			}
		}

		err = s.chatRepo.Create(ctx, &models.Chat{
			Role: "assistant",
			Message: fullSentences,
			SessionRoomId: sessionId,
			ReceivedTime: time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		

		if err != nil {
			return apperr.Internal(fmt.Errorf("chat service chat stream on finish: %w", err))
		}
		return onFinish(fullSentences)
	})
}

func (s *chatService) FindRecentMessages(ctx context.Context, sessionId string, numLastMessages int32) ([]*models.Chat, error) {
	return s.chatRepo.FindRecentMessages(ctx, sessionId, numLastMessages)
}