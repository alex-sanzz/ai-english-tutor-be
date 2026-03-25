package postgres

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"ai-tutor-backend/internal/repository"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type subcriptionRepository struct {
	db     *pgxpool.Pool
	logger log.Logger
}

func NewSubcriptionRepository(db *pgxpool.Pool,logger log.Logger) repository.SubcriptionRepository{
	return subcriptionRepository{
		db: db,
		logger: logger,
	}
}

// CREATE TABLE subcriptions(
//     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     user_id TEXT NOT NULL,
//     purchase_token TEXT NOT NULL,
//     purchase_time TIMESTAMPTZ NOT NULL,
//     expiry_at TIMESTAMPTZ NOT NULL,
//     auto_renewing BOOLEAN NOT NULL,
//     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
// );


func (r subcriptionRepository) Create(ctx context.Context, p models.Subcription) error{
	query := "INSERT INTO subcriptions(user_id, purchase_token, purchase_time, expiry_at, auto_renewing) VALUES($1, $2, $3, $4, $5)"

	_, err := r.db.Query(ctx, query, p.UserId, p.PurchaseToken, p.PurchaseTime, p.ExpiryAt, p.AutoRenewing)

	if err != nil {
		return apperr.Internal(fmt.Errorf("purchase history repository create error: %w", err))
	}

	return nil
}

func (r subcriptionRepository) FindActiveSubscription(ctx context.Context, userId string) (*models.Subcription, error){
	query := "SELECT id, user_id, purchase_token, purchase_time, expiry_at, auto_renewing FROM subcriptions WHERE user_id = $1 AND expiry_at >= NOW() LIMIT 1"

	var s models.Subcription

	err := r.db.QueryRow(ctx, query, userId).Scan(&s.ID, &s.UserId, &s.PurchaseToken, &s.PurchaseTime, &s.ExpiryAt, &s.AutoRenewing)

	if err != nil {
		return nil, apperr.Internal(fmt.Errorf("subcription repository find active subcription error: %w", err))
	}

	return &s, nil  

}