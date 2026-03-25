package router

import (
	"ai-tutor-backend/internal/handler"
	"ai-tutor-backend/internal/infrastructure/jwt"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/middleware"
	"time"

	// "github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func SetupRoute(c *gin.Engine, messageHandler *handler.MessageHandler, authHandler *handler.AuthHandler, sessionRoomHandler *handler.SessionRoomHandler, conversationQuestionHandler *handler.ConversationQuestionHandler, paymentHandler *handler.PaymentHandler, jwtClient jwt.JwtClient, logger log.Logger){

	apiGroup := c.Group("/api")

	// apiGroup.Use(timeout.New(timeout.WithTimeout(15*time.Second)))

	protectedGroup := apiGroup.Group("")
	
	protectedGroup.Use(middleware.JwtMiddleware(jwtClient, logger))

	protectedGroup.GET("/rooms", sessionRoomHandler.FindAllTopics)
	protectedGroup.POST("/rooms", sessionRoomHandler.CreateSessionRoom)
	protectedGroup.DELETE("/rooms/:id", sessionRoomHandler.DeleteSessionRoom)

	protectedGroup.GET("/questions", conversationQuestionHandler.FindAll)
	protectedGroup.GET("/questions/:id", conversationQuestionHandler.FindById)
	protectedGroup.POST("/questions", conversationQuestionHandler.GenerateQuestion)
	protectedGroup.POST("/questions/answer", conversationQuestionHandler.AnswerQuestion)

	protectedGroup.POST("/payment/subscription/ack", paymentHandler.AcknowledgeSubcription)

	sseGroup := protectedGroup.Group("/sse")
	sseGroup.GET("", messageHandler.RegisterSse)
	sseGroup.POST("/send/:id", messageHandler.SendMessage)
	sseGroup.POST("/send/:id/audio", messageHandler.SendAudioMessage)
	authGroup := apiGroup.Group("/auth")

	authGroup.POST("/google/login", authHandler.GoogleLogin)

	authGroup.POST("/refresh", authHandler.RefreshToken)

	apiGroup.GET("/:sessionId/messages", messageHandler.FindRecentMessages)
	
	protectedGroup.GET("/health", func(c *gin.Context){
		c.JSON(200, gin.H{
			"status": "ok",
			"timestamp": time.Now().String(),
			"service": "ai-tutor backend service",
		})
	})

	


}