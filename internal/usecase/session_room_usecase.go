package usecase

import (
	"ai-tutor-backend/internal/models"
	"ai-tutor-backend/internal/service"
	"context"
)

type SessionRoomUseCase struct {
	sessionRoomService service.SessionRoomService
}

func NewSessionRoomUseCase(sessionRoomService service.SessionRoomService) *SessionRoomUseCase {
	return &SessionRoomUseCase{
		sessionRoomService: sessionRoomService,
	}
}

func (u *SessionRoomUseCase) CreateSessionRoom(ctx context.Context,userId string, roomType string, icon string, topic string) (string, error) {
	return u.sessionRoomService.CreateSessionRoom(ctx, userId, roomType, icon, topic)
}

func (u *SessionRoomUseCase) FindAllRooms(ctx context.Context,userId string, roomType string) ([]*models.SessionRoom, error){
	return u.sessionRoomService.FindAll(ctx, userId, roomType)
}

func (u *SessionRoomUseCase) FindById(ctx context.Context, id string) (*models.SessionRoom, error){
	return u.sessionRoomService.FindById(ctx, id)
}

func (u *SessionRoomUseCase) DeleteById(ctx context.Context, id string) error{
	return u.sessionRoomService.DeleteByID(ctx, id)
}