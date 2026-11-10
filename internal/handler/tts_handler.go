package handler

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/dto"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/usecase"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TtsHandler struct {
	ttsUsecase usecase.TTSUsecase
	logger log.Logger
}

func (t TtsHandler) initiateLogger(c context.Context) (log.Logger){
	return t.logger.FromContext(c)
}


func NewTtsHandler(ttsUsecase usecase.TTSUsecase, logger log.Logger) *TtsHandler{
	return &TtsHandler{
		ttsUsecase: ttsUsecase,
		logger: logger,
	}
}

func (t TtsHandler) Generate(c *gin.Context){
	logger := t.initiateLogger(c.Request.Context())

	var reqBody dto.TTSRequestDto

	if err := c.BindJSON(&reqBody); err != nil {
		logger.Error("tts handler parsing request body to struct error", zap.Error(err))
		writeError(c, apperr.BadRequest("400", "invalid request body format", fmt.Errorf("tts handler parsing request body to struct error: %w", err)))
		return
	}

	audioSource, err := t.ttsUsecase.Generate(c.Request.Context(), reqBody.Text)

	if err != nil {
		logger.Error("tts handler generate method error", zap.Error(err))
		writeError(c, apperr.Internal(err))
		return
	}

	defer audioSource.Close()

	c.Header("Content-Type", "audio/pcm")
	c.Header("Cache-Control", "no-store")
	c.Status(http.StatusOK)
	
	buffer := make([]byte, 4096)

	for {
		bytes, readErr := audioSource.Read(buffer)

		if bytes > 0 {
			if _, err := c.Writer.Write(buffer[:bytes]); err != nil {
			    logger.Error("tts handler generate method error", zap.Error(err))
				writeError(c, readErr)
				return
			}

			c.Writer.Flush()
		}

		if readErr == io.EOF {
			break
		}


		if readErr != nil {
			logger.Error("tts handler generate method error", zap.Error(err))
			writeError(c, readErr)
			return
		}

	}

	

	
}