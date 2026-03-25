package repository

import (
	"ai-tutor-backend/internal/models"
	"context"
)

type ConversationQuestionRepository interface {
	FindById(ctx context.Context, id string) (*models.ConversationQuestion, error)
	FindAll(ctx context.Context, sessionRoomId string, limit *int, desc bool) ([]*models.ConversationQuestion, error)
	Create(ctx context.Context, m *models.ConversationQuestion) (string, error)
	CreateBatch(ctx context.Context, questions []*models.ConversationQuestion) error
	CheckIfThereIsAnUnansweredQuestion(ctx context.Context, sessionRoomId string) (bool, error)
	FindOneUnansweredQuestion(ctx context.Context, sessionRoomId string) (*models.ConversationQuestion, error)
	UpdateAnswer(ctx context.Context, id, answer, reviewResult string) error
	Delete(ctx context.Context, id string) error
}