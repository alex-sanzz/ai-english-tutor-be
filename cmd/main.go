package main

import (
	"ai-tutor-backend/internal/config"
	"ai-tutor-backend/internal/handler"
	"ai-tutor-backend/internal/infrastructure/assemblyai"
	androidpublisher "ai-tutor-backend/internal/infrastructure/google/android-publisher"
	"ai-tutor-backend/internal/infrastructure/google/auth"
	"ai-tutor-backend/internal/infrastructure/httpclient"
	"ai-tutor-backend/internal/infrastructure/jwt"
	"ai-tutor-backend/internal/infrastructure/openai"
	"ai-tutor-backend/internal/infrastructure/postgres"
	"ai-tutor-backend/internal/infrastructure/sse"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/middleware"
	"ai-tutor-backend/internal/seed"

	"ai-tutor-backend/internal/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "ai-tutor-backend/internal/infrastructure/gemini"

	"ai-tutor-backend/internal/router"
	"ai-tutor-backend/internal/usecase"

	// "context"

	// "fmt"

	"github.com/gin-gonic/gin"
	openaigo "github.com/openai/openai-go/v3"
	"go.uber.org/zap"
	// "google.golang.org/genai"
)

func main() {

	logger, _ := log.NewZapLogger("debug", false)

	defer logger.Sync()

	cfg, err := config.InitConfig()

	if err != nil {

		logger.Fatal("failed init config", zap.Error(err))
	}
	logger.Info("config loaded successfully")
	logger.Info(fmt.Sprintf("database name: %s", cfg.Database.DatabaseName))
	dbDsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, fmt.Sprint(cfg.Database.Port), cfg.Database.DatabaseName)

	if err := postgres.RunMigrations(dbDsn); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	g := gin.Default()
	g.Use(middleware.RecoveryMiddleware(logger))
	// ctx := context.Background()
	sseBroker := sse.NewSseBroker(logger)
	// aiClient, err := genai.NewClient(ctx, nil)

	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }
	// geminiClient := gemini.NewGeminiClient(aiClient)

	// geminiUseCase := usecase.NewGeminiUseCase(geminiClient)
	pgPool, err := postgres.NewPool(context.Background(), dbDsn)

	seeder := seed.NewSeeder(pgPool)

	if err != nil {
		logger.Fatal("failed init postgresql", zap.Error(err))
	}

	err = seeder.SeedRoomType(context.Background())

	if err != nil {
		logger.Fatal("failed seed room type", zap.Error(err))
	}
	
	logger.Info("connected to postgresql successfully, host: " + cfg.Database.Host + ", port: " + fmt.Sprint(cfg.Database.Port))
	chatRepository := postgres.NewChatRepository(pgPool, logger)
	openAiInstance := openaigo.NewClient()
	httpConfig := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}
	httpClient := httpclient.NewHttpClient(httpConfig)
	assemblyAiClient := assemblyai.NewAssemblyAiClient(cfg.AssemblyAi, httpClient, logger)
	transcriptionService := service.NewTranscribeService(assemblyAiClient)
	openAiClient := openai.NewOpenAiClient(openAiInstance, cfg, logger)
	chatService := service.NewChatService(openAiClient, chatRepository, logger, cfg.Ai)
	openAiUseCase := usecase.NewOpenAiUseCase(chatService, logger)

	sessionRoomRepository := postgres.NewSessionRoomRepository(pgPool, logger)
	sessionRoomService := service.NewSessionRoomService(sessionRoomRepository, chatRepository)
	sessionRoomUseCase := usecase.NewSessionRoomUseCase(sessionRoomService)
	sessionRoomHandler := handler.NewSessionRoomHandler(sessionRoomUseCase, logger)
	transcribeUseCase := usecase.NewTranscribeUseCase(transcriptionService)
	messageHandler := handler.NewMessageHandler(sseBroker, openAiUseCase, sessionRoomUseCase, transcribeUseCase, logger, &cfg.Ai)
	
	googleAuthClient := auth.NewGoogleAuthCLient(&cfg.GoogleAuth, logger)
	userRepo := postgres.NewUserRepository(pgPool)
	jwtClient, err := jwt.NewRsaJwtClient(cfg.Jwt, logger)
	if err != nil {
		logger.Fatal("failed to init jwt client", zap.Error(err))

	}

	refreshTokenRepo := postgres.NewRefreshTokenRepository(pgPool, logger)

	authService := service.NewAuthService(userRepo, googleAuthClient, jwtClient, refreshTokenRepo, &cfg.Jwt, logger)
	googleAuthUseCase := usecase.NewAuthUseCase(authService)
	authHandler := handler.NewAuthHandler(googleAuthUseCase, logger)
	
	conversationQuestionRepo := postgres.NewConversationQuestionRepository(pgPool, logger)
	conversationQuestionSvc := service.NewConversationQuestionService(conversationQuestionRepo, sessionRoomRepository, openAiClient, assemblyAiClient, logger,  cfg.Ai)
	conversationQuestionUseCase := usecase.NewConversationQuestionUseCase(conversationQuestionSvc, logger)
	conversationQuestionHandler := handler.NewConversationQuestionHandler(*conversationQuestionUseCase, logger)

	androidPublisherClient, err := androidpublisher.New(context.Background(), "./files/google-account-service.json")
	
	if err != nil {
		logger.Fatal("failed to init android publisher", zap.Error(err))

	}
	
	subcriptionRepo := postgres.NewSubcriptionRepository(pgPool, logger)
	paymentSvc := service.NewPaymentService(androidPublisherClient, userRepo, subcriptionRepo)
	paymentUseCase := usecase.NewPaymentUseCase(paymentSvc)
	paymentHandler := handler.NewPaymentHandler(*paymentUseCase)

	router.SetupRoute(g, messageHandler, authHandler, sessionRoomHandler, conversationQuestionHandler, paymentHandler, jwtClient, logger)

	go sseBroker.Listen()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: g,
		// Maxium time the client can send the full request
		ReadTimeout: 15 * time.Second,
		// // Maximum time the server can respond to the client
		// WriteTimeout: 15 * time.Second,
		// Maximum time to keep an idle keep-alive connection open.
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		logger.Info("server is running port :" + fmt.Sprint(cfg.Port))
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal("failed to listen and serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	// When SIGTERM signal is received, for example: when you pressing ctrl + c
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// remember! channel will wait until it receives a value
	<-quit

	logger.Info("syscall.SIGINT or syscall.SIGTERM is received, shutting down server...")
	// this creates a custom context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// shutdown means stop accepting new connections and finish ongoing request (or don't interupt any ongoing request)
	// But if both of them are taking longer than 10 second to finish, then just ignore it immediately or just continue the process
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server is forced to shutdown: %v", zap.Error(err))
	}

	logger.Info("server is off gracefully")
}
