package service

import (
	"ai-tutor-backend/internal/models"
	"context"
)

type SessionRoomService interface {
	CreateSessionRoom(context context.Context, userId, roomType, icon, topic string) (string, error)
	FindAll(context context.Context, userId, roomType string) ([]*models.SessionRoom, error)
	FindById(ctx context.Context, id string) (*models.SessionRoom, error)
	DeleteByID(ctx context.Context, id string) error 
}