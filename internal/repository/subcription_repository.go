package repository

import (
	"ai-tutor-backend/internal/models"
	"context"
)

type SubcriptionRepository interface {
	Create(ctx context.Context, p models.Subcription) error
	FindActiveSubscription(ctx context.Context, userId string) (*models.Subcription, error)
}
