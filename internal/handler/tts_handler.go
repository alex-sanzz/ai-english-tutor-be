package handler

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/dto"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/usecase"
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

func NewTtsHandler(ttsUsecase usecase.TTSUsecase, logger log.Logger) *TtsHandler{
	return &TtsHandler{
		ttsUsecase: ttsUsecase,
		logger: logger,
	}
}

func (t TtsHandler) Generate(c *gin.Context){
	var reqBody dto.TTSRequestDto

	if err := c.BindJSON(&reqBody); err != nil {
		t.logger.Error("tts handler parsing request body to struct error", zap.Error(err))
		writeError(c, apperr.BadRequest("400", "invalid request body format", fmt.Errorf("tts handler parsing request body to struct error: %w", err)))
		return
	}

	audioSource, err := t.ttsUsecase.Generate(c.Request.Context(), reqBody.Text)

	if err != nil {
		t.logger.Error("tts handler generate method error", zap.Error(err))
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
				t.logger.Error("tts handler generate method error", zap.Error(err))
				writeError(c, readErr)
				return
			}

			c.Writer.Flush()
		}

		if readErr == io.EOF {
			break
		}


		if readErr != nil {
			t.logger.Error("tts handler generate method error", zap.Error(err))
			writeError(c, readErr)
			return
		}

	}

	
}