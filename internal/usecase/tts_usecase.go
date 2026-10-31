package usecase

import (
	"ai-tutor-backend/internal/service"
	"context"
	"io"
)

type TTSUsecase struct {
	ttsService service.TTSService
}

func NewTTSUsecase(ttsService service.TTSService) *TTSUsecase{
	return &TTSUsecase{
		ttsService: ttsService,
	}
}

func (t *TTSUsecase) Generate(ctx context.Context, text string) (io.ReadCloser, error){
	return t.ttsService.Generate(ctx, text)
}