package service

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/config"
	"ai-tutor-backend/internal/infrastructure/google/auth"
	"ai-tutor-backend/internal/infrastructure/jwt"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"ai-tutor-backend/internal/repository"
	"context"
	"fmt"
	"time"
)

type authService struct {
	jwtClient jwt.JwtClient
	userRepo repository.UserRepository
	googleAuthClient auth.GoogleAuthClient
	refreshTokenRepository repository.RefreshTokenRepository
	jwtCfg *config.JwtConfig
	logger log.Logger 
}

func NewAuthService(userRepo repository.UserRepository, googleAuthClient auth.GoogleAuthClient, jwtClient jwt.JwtClient, refreshTokenRepository repository.RefreshTokenRepository, cfg *config.JwtConfig, logger log.Logger) AuthService{
	return authService{
		userRepo: userRepo,
		jwtClient: jwtClient,
		googleAuthClient: googleAuthClient,
		refreshTokenRepository: refreshTokenRepository,
		jwtCfg: cfg,
		logger: logger,
	}
}

func (a authService) GoogleLogin(ctx context.Context, googleIdToken string, nonce string) (string, string, error){
	payload, err := a.googleAuthClient.ParseAndVerifyGoogleTokenId(ctx, googleIdToken, nonce)

	if err != nil {
		return "", "", apperr.Unauthorized("unauthorized", "unauthorized", fmt.Errorf("auth service implementation google login error: %w", err))
	}

	userId, err := a.userRepo.Create(ctx, models.User{
		GoogleSub: payload.Subject,
		Name: payload.Claims["name"].(string),
		Email: payload.Claims["email"].(string),
	})

	if err != nil {
		return "", "", apperr.Internal(fmt.Errorf("auth service implementation google login error: %w", err))
	}

	accessToken, refreshToken, err := a.GenerateTokenPair(ctx, userId)

	if err != nil {
		return "", "", apperr.Internal(fmt.Errorf("auth service implementation google login error: %w", err))
	}
	

	return accessToken, refreshToken, nil 

}

func (a authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error){
	refreshTokenModel, err := a.refreshTokenRepository.FindByToken(ctx, refreshToken)

	if err != nil {
		return "", "", apperr.Internal(fmt.Errorf("auth service refresh token error: %w", err))
	}

	if refreshTokenModel == nil {
		return "", "", apperr.BadRequest("unauthorized", "unauthorized", fmt.Errorf("auth service refresh token error: requested refresh token not found"))
	}

	err = a.refreshTokenRepository.RevokeByToken(ctx, refreshToken)

	if err != nil {
		return "", "", apperr.Internal(fmt.Errorf("auth service refresh token error: requested refresh token not found"))
	}

	return a.GenerateTokenPair(ctx, refreshTokenModel.UserId)


}

func (a authService) GenerateTokenPair(ctx context.Context, userId string) (string, string, error){
	accessToken, err := a.jwtClient.GenerateAccessToken(userId)

	if err != nil {
		return "", "", apperr.Internal(fmt.Errorf("auth service generate token pair error: %w", err))
	}

	refreshToken, err := a.jwtClient.GenerateRefreshToken()

	if err != nil {
		return "", "", apperr.Internal(fmt.Errorf("auth service generate token pair error: %w", err))
	}

	now := time.Now()

	err = a.refreshTokenRepository.Create(ctx, &models.RefreshToken{
		UserId: userId,
		Token: refreshToken,
		ExpiresAt: now.Add(a.jwtCfg.RefreshTokenTTL),
		CreatedAt: now,
	})

	if err != nil {
		return "", "", apperr.Internal(fmt.Errorf("auth service generate token pair error: %w", err))
	}

	return accessToken, refreshToken, nil 
}