package service

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/infrastructure/openai"
	"ai-tutor-backend/internal/log"
	"context"
	"fmt"
	"io"
)

type ttsService struct {
	openaiclient *openai.OpenAiClient
	logger log.Logger 
}

func NewTtsService(openaiclient *openai.OpenAiClient, logger log.Logger ) TTSService {
	return ttsService{
		openaiclient: openaiclient,
		logger: logger,
	}
}

func (t ttsService) Generate(ctx context.Context, text string) (io.ReadCloser, error){
	s, err := t.openaiclient.StreamTTS(ctx, text)

	if err != nil {
		return nil, apperr.Internal(fmt.Errorf("tts service impl generate error: %w", err))
	}

	return s, nil
}

