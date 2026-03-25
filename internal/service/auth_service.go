package service

import "context"

type AuthService interface {
	GoogleLogin(ctx context.Context, googleIdToken string, nonce string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	GenerateTokenPair(ctx context.Context, userId string) (string, string, error)
}