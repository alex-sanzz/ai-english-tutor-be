package usecase

import (
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"ai-tutor-backend/internal/service"
	"context"
)

type ConversationQuestionUseCase struct {
	service service.ConversationQuestionService
	logger log.Logger
}

func NewConversationQuestionUseCase(service service.ConversationQuestionService, logger log.Logger) *ConversationQuestionUseCase{
	return &ConversationQuestionUseCase{
		service: service,
		logger: logger,
	}
}

func (u ConversationQuestionUseCase) FindById(ctx context.Context, id string) (*models.ConversationQuestion, error){

	return u.service.FindById(ctx, id)
}

func (u ConversationQuestionUseCase) FindAll(ctx context.Context, sessionRoomId string) ([]*models.ConversationQuestion, error){
	limit := 50
	return u.service.FindAll(ctx, sessionRoomId, &limit, false)
}

func (u ConversationQuestionUseCase) FindAllAnsweredQuestion(ctx context.Context, sessionRoomId string) ([]*models.ConversationQuestion, error){
	return u.service.FindAllAnsweredQuestion(ctx, sessionRoomId)
}

func (u ConversationQuestionUseCase) GenerateQuestion(ctx context.Context, sessionRoomId string) error{
	return u.service.GenerateQuestion(ctx, sessionRoomId)
}

func (u ConversationQuestionUseCase) AnswerQuestion(ctx context.Context, id string, alternateVersion string, culturalContext string, paraphraseVersion string, answer []byte) error{
	return u.service.AnswerQuestion(ctx, id, alternateVersion, culturalContext, paraphraseVersion, answer)
}