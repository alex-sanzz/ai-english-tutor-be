package handler

import (
	"ai-tutor-backend/internal/dto"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
	logger log.Logger
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase, logger log.Logger) *AuthHandler{
	return &AuthHandler{
		authUseCase: authUseCase,
		logger: logger,
	}
}

func (a *AuthHandler) GoogleLogin(c *gin.Context){
	var request dto.GoogleTokenIdVerifyRequest

	if err := c.BindJSON(&request); err != nil {
		a.logger.Error("auth handler google login error: ", zap.Error(err))
		writeError(c, err)
		return
	}

	accessToken, refreshToken, err := a.authUseCase.GoogleLogin(c.Request.Context(), request.GoogleTokenId, request.Nonce)
	
	if err != nil {
		a.logger.Error("auth handler google login error: ", zap.Error(err))
		writeError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}


func (a *AuthHandler) RefreshToken(c *gin.Context){
	var request dto.RefreshTokenRequestDto

	if err := c.BindJSON(&request); err != nil {
		a.logger.Error("auth handler bind json error", zap.Error(err))
		writeError(c, err)
		return
	}

	accessToken, refreshToken, err := a.authUseCase.RefreshToken(c.Request.Context(), request.RefreshToken)
	
	if err != nil {
		a.logger.Error("auth handler refresh token service error", zap.Error(err))
		writeError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}
