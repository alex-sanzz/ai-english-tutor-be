package middleware

import (
	"ai-tutor-backend/internal/log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context){

		defer func(){
			if err := recover(); err != nil {
				logger.Error("recovery middleware error: ", zap.Any("error", err), 
				zap.String("stack", string(debug.Stack())),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))
				c.JSON(500, gin.H{
					"error_message": "internal server error",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}