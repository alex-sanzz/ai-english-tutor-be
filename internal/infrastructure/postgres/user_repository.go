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

type UserRepository struct {
	db *pgxpool.Pool
	logger log.Logger
}

func NewUserRepository(db *pgxpool.Pool) (repository.UserRepository){
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(ctx context.Context, user models.User) (string, error){
	query := "INSERT INTO users (google_sub, email, name) VALUES($1, $2, $3) ON CONFLICT (google_sub) DO UPDATE SET last_login_at = NOW() RETURNING id"
	
	var id string 
	
	err := u.db.QueryRow(ctx, query, user.GoogleSub, user.Email, user.Name).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("user repository create error: %w", err) 
	}

	return id, nil
} 

func (u *UserRepository) FindById(ctx context.Context, id string) (*models.User, error){
	query := "SELECT id, google_sub, email, name, created_at, last_login_at FROM users WHERE id = $1 LIMIT 1"

	var user models.User

	err := u.db.QueryRow(ctx, query).Scan(&user.ID, &user.GoogleSub, &user.Email, &user.Name, &user.CreatedAt, &user.LastLoginAt)

	if err != nil {
		return nil, apperr.Internal(fmt.Errorf("user repository find by id error: %w", err))
	}

	return &user, nil
}

