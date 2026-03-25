package middleware

import (
	"ai-tutor-backend/internal/infrastructure/jwt"
	"ai-tutor-backend/internal/log"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func JwtMiddleware(jwtClient jwt.JwtClient, logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context){
	
		authHeader := c.GetHeader("Authorization")
	
		if authHeader == "" {
			logger.Warn("jwt middleware warning: there is no authorization header")
			c.JSON(401, gin.H{"error_message": "unauthorized"})
			c.Abort()
			return 
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("jwt middlware warning: invalid authorization format", zap.String("authorization value", authHeader))
			c.JSON(401, gin.H{"error_message": "unauthorized"})
			c.Abort()
			return 
		}

		token := parts[1]

		claims, err := jwtClient.ParseAndValidate(token)

		if err != nil {
			logger.Error("jwt middlware error: ", zap.Error(err))
			c.JSON(401, gin.H{"error_message": "unauthorized"})
			c.Abort()
			return
		}

		c.Set("userId", claims.Subject)

		c.Next()
		
	} 

}