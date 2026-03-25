package repository

import (
	"ai-tutor-backend/internal/models"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (string, error)
	FindById(ctx context.Context, id string) (*models.User, error)
}