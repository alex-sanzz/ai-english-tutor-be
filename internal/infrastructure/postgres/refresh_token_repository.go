package postgres

import (
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type refreshTokenRepository struct {
	db *pgxpool.Pool
	logger log.Logger
}

func NewRefreshTokenRepository(db *pgxpool.Pool, logger log.Logger) *refreshTokenRepository{
	return &refreshTokenRepository{
		db: db,
		logger: logger,
	}
}

// Create(ctx context.Context, refreshToken *models.RefreshToken) error 
// 	FindByToken(ctx context.Context, token string) (*models.RefreshToken, error)
// 	RevokeByToken(ctx context.Context, token string) error 
// 	RevokeAllByUserId(ctx context.Context, userId string) error 
// 	DeleteExpired(ctx context.Context) error

func (r *refreshTokenRepository) Create(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := "INSERT INTO refresh_tokens(user_id, token, expires_at) VALUES ($1, $2, $3)"

	_, err := r.db.Exec(ctx, query, refreshToken.UserId, refreshToken.Token, refreshToken.ExpiresAt)

	if err != nil {
		return fmt.Errorf("refresh token repository create error: %w", err)
	}

	return nil 
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error){
	query := "SELECT * FROM refresh_tokens WHERE token = $1 LIMIT 1"

	var refreshToken models.RefreshToken

	err := r.db.QueryRow(ctx, query, token).Scan(&refreshToken.ID, &refreshToken.UserId, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.CreatedAt, &refreshToken.RevokedAt)


	if err == sql.ErrNoRows {
		return nil, nil 
	}

	if err != nil {
		return nil, fmt.Errorf("refresh token repository find by token error: %w", err)
	}

	return &refreshToken, nil 
}

func (r *refreshTokenRepository) RevokeByToken(ctx context.Context, token string) error {
	query := "UPDATE refresh_tokens SET revoked_at = NOW() WHERE token = $1 AND revoked_at = NULL"

	_, err := r.db.Exec(ctx, query, token)

	if err != nil {
		return fmt.Errorf("refresh token repository revoke by token error: %w", err)
	}

	return nil 
}

func (r *refreshTokenRepository) RevokeAllByUserId(ctx context.Context, userId string) error {
	query := "UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at = NULL"

	_, err := r.db.Exec(ctx, query, userId)

	if err != nil {
		return fmt.Errorf("refresh token repository revoke all by user id error: %w", err)
	}

	return nil 
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := "DELETE FROM refresh_tokens WHERE revoked_at IS NOT NULL OR expires_at <= NOW() "

	_, err := r.db.Exec(ctx, query)

	if err != nil {
		return fmt.Errorf("refresh token repository delete expired error: %w", err)
	}

	return nil 
}
