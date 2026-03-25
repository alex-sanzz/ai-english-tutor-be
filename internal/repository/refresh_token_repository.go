package repository

import (
	"ai-tutor-backend/internal/models"
	"context"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *models.RefreshToken) error 
	FindByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeByToken(ctx context.Context, token string) error 
	RevokeAllByUserId(ctx context.Context, userId string) error 
	DeleteExpired(ctx context.Context) error
}