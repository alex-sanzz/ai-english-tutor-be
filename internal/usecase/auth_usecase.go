package usecase

import (

	"ai-tutor-backend/internal/service"
	"context"
	"fmt"

)

type AuthUseCase struct {
	authService service.AuthService
}

func NewAuthUseCase(authService service.AuthService) *AuthUseCase {
	return &AuthUseCase{
		authService: authService,
	}
}

func (g *AuthUseCase) GoogleLogin(context context.Context, googleIdToken string, nonce string) (string, string, error) {
	accessToken, refreshToken, err := g.authService.GoogleLogin(context, googleIdToken, nonce)

	if err != nil {
		return "", "", fmt.Errorf("auth use case google login error : %w", err)
	}

	return accessToken, refreshToken, nil

	
}

func (g *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	accessToken, refreshToken, err := g.authService.RefreshToken(ctx, refreshToken)

	if err != nil {
		return "", "", fmt.Errorf("auth use case refresh token error: %w", err)
	}

	return accessToken, refreshToken, err
}