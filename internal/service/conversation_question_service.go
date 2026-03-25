package service

import (
	"ai-tutor-backend/internal/models"
	"context"
)

type ConversationQuestionService interface {
	FindById(ctx context.Context, id string) (*models.ConversationQuestion, error)
	FindAll(ctx context.Context, sessionRoomId string, limit *int, desc bool) ([]*models.ConversationQuestion, error)
	GenerateQuestion(ctx context.Context, sessionRoomId string) error
	AnswerQuestion(ctx context.Context, id string, answer []byte) error
}