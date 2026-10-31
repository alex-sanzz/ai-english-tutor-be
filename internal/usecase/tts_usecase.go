package usecase

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/service"
	"context"
	"fmt"
	"io"
)

type TTSUsecase struct {
	ttsService service.TTSService
	logger log.Logger
}

func NewTTSUsecase(ttsService service.TTSService, logger log.Logger) *TTSUsecase{
	return &TTSUsecase{
		ttsService: ttsService,
		logger: logger,
	}
}

func (t *TTSUsecase) Generate(ctx context.Context, text string) (io.ReadCloser, error){
	c, err := t.ttsService.Generate(ctx, text)
	if err != nil {
		return nil, apperr.Internal(fmt.Errorf("tts usecase generate method error: %w", err))
	}
	return c, nil
}