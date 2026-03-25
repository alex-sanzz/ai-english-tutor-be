package repository

import (
	"ai-tutor-backend/internal/models"
	"context"
)

type SessionRoomRepository interface {
	CreateSessionRoom(ctx context.Context, userId, roomType, icon, topic string) (string, error) 
	FindAllByUserIdAndRoomType(ctx context.Context, userId, roomType string) ([]*models.SessionRoom, error)
	FindById(ctx context.Context, id string) (*models.SessionRoom, error)
	DeleteById(ctx context.Context, id string) error 
}