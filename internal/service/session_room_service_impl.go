package service

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/models"
	"ai-tutor-backend/internal/repository"
	"context"
	"fmt"
)

type sessionRoomService struct {
	sessionRoomRepo repository.SessionRoomRepository
}

func NewSessionRoomService(sessionRoomRepo repository.SessionRoomRepository) *sessionRoomService {
	return &sessionRoomService{
		sessionRoomRepo: sessionRoomRepo,
	}
}

func (s *sessionRoomService) CreateSessionRoom(context context.Context, userId, roomType, icon, topic string) (string, error){
	sessionRoomId, err := s.sessionRoomRepo.CreateSessionRoom(context, userId, roomType, icon, topic)
	if err != nil {
		return "", apperr.Internal(fmt.Errorf("session room service create session room error: %w", err))
	}
	return sessionRoomId, nil
}

func (s *sessionRoomService) FindAll(context context.Context, userId, roomType string) ([]*models.SessionRoom, error){
	return s.sessionRoomRepo.FindAllByUserIdAndRoomType(context, userId, roomType)
}

func (s *sessionRoomService) FindById(ctx context.Context, id string) (*models.SessionRoom, error) {
	return  s.sessionRoomRepo.FindById(ctx, id)
}

func (s *sessionRoomService) DeleteByID(ctx context.Context, id string) error {
	return s.sessionRoomRepo.DeleteById(ctx, id)
}