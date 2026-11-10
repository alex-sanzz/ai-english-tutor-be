package middleware

import (
	"ai-tutor-backend/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func SetRequestIdMiddleware(base log.Logger) gin.HandlerFunc {
	return func(c *gin.Context){
		id := c.GetHeader("X-Request-ID")

		if id == "" {
			id = uuid.NewString()
		}

		c.Header("X-Request-ID", id)

		// Set additional log field for every next log
		child := base.With(zap.String("request_id", id))

		// create a new context, with a new logger
		ctx := log.ContextWithLogger(c.Request.Context(), child)

		c.Request = c.Request.WithContext(ctx)

		c.Next()

	}
}