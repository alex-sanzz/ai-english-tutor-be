package service

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/config"
	"ai-tutor-backend/internal/infrastructure"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"ai-tutor-backend/internal/repository"
	"context"
	"fmt"

	"strings"
	"time"
)

type conversationQuestionService struct {
	repo repository.ConversationQuestionRepository
	sessionRoomRepo repository.SessionRoomRepository
	aiClient infrastructure.AiChatClient
	transcriptionClient infrastructure.TranscriptionClient
	cfg config.AiConfig
	logger log.Logger
	
}

func NewConversationQuestionService(repo repository.ConversationQuestionRepository, sessionRoomRepo repository.SessionRoomRepository, aiClient infrastructure.AiChatClient, transcriptionClient infrastructure.TranscriptionClient, logger log.Logger, cfg config.AiConfig) ConversationQuestionService{
	return conversationQuestionService{
		repo: repo, aiClient: aiClient, sessionRoomRepo: sessionRoomRepo, cfg: cfg, transcriptionClient: transcriptionClient, logger: logger,
	}
}

func (s conversationQuestionService) FindById(ctx context.Context, id string) (*models.ConversationQuestion, error){
	return s.repo.FindById(ctx, id)
}

func (s conversationQuestionService) FindAll(ctx context.Context, sessionRoomId string, limit *int, desc bool) ([]*models.ConversationQuestion, error){
	questions, err := s.repo.FindAll(ctx, sessionRoomId, limit, desc)

	if err != nil {
		return nil, err 
	}

	if len(questions) == 0 {
		err = s.GenerateQuestion(ctx, sessionRoomId)
		if err != nil {
			return nil, err 
		}
		questions, err = s.FindAll(ctx, sessionRoomId, limit, desc)
		if err != nil {
			return nil, err 
		}
	}

	return questions, nil 
}

func (s conversationQuestionService) GenerateQuestion(ctx context.Context, sessionRoomId string) error{
	
	room, err := s.sessionRoomRepo.FindById(ctx, sessionRoomId)

	if err != nil {
		return err
	}

	q, err := s.repo.FindOneUnansweredQuestion(ctx, sessionRoomId)

	if err != nil {
		return fmt.Errorf("conversation question service impl error: %w", err)
	}

	if q != nil {
		return nil 
	}

	answer, err := s.aiClient.AskQuestion(ctx, strings.ReplaceAll(s.cfg.GenerateQuestionPrompt, "{{topic}}", room.Topic))

	if err != nil {
		return fmt.Errorf("conversation question service impl ai client ask question error: %w", err)
	}


	answers := strings.Split(answer, ",")

	questionModels := []*models.ConversationQuestion{}

	

	for _, a := range answers {
		questionModels = append(questionModels, &models.ConversationQuestion{
			Question: a,
			SessionRoomId: sessionRoomId,
			CreatedAt: time.Now(),
		})
		time.Sleep(10 *time.Millisecond)
	}

	err = s.repo.CreateBatch(ctx, questionModels)

	if err != nil {
		return fmt.Errorf("conversation question service impl create batch error: %w", err)
	}


	return nil 

}

func (s conversationQuestionService) AnswerQuestion(ctx context.Context, id string, answer []byte) error{
	data, err := s.repo.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("onversation question service answer question error: %w", err)
	}
	uploadedUrl, err := s.transcriptionClient.UploadAudio(ctx, answer)

	if err != nil {
		return apperr.Internal(fmt.Errorf("conversation question service answer question error: can't interpret the message")) 
	}

	transcribedText, err := s.transcriptionClient.TranscribeAudio(ctx, uploadedUrl)
	
	if err != nil {
		return apperr.Internal(fmt.Errorf("conversation question service answer question error: can't interpret the message")) 
	}

	if transcribedText == ""{
		return apperr.BadRequest("400", "can't interpret the message, please record the audio again or check the mic", fmt.Errorf("conversation question service answer question error: can't interpret the message"))
	}

	reviewAnswer, err := s.aiClient.AskQuestion(ctx, strings.ReplaceAll(strings.ReplaceAll(s.cfg.OneSentenceEnglishEvaluation, "{{user_answer}}", transcribedText), "{{question}}", data.Question))
	
	if err != nil {
		return err
	}

	err = s.repo.UpdateAnswer(ctx, id, transcribedText, reviewAnswer)

	if err != nil {
		return err 
	}

	return nil

}
