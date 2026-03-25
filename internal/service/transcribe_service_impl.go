package service

import (
	"ai-tutor-backend/internal/infrastructure"
	"context"
)

type transcribeService struct {
	transcriptionClient infrastructure.TranscriptionClient
}

func NewTranscribeService(transcriptionClient infrastructure.TranscriptionClient) *transcribeService {
	return &transcribeService{
		transcriptionClient: transcriptionClient,
	}
}

func (s *transcribeService) UploadAudioFile(ctx context.Context, fileBytes []byte) (string, error) {
	return s.transcriptionClient.UploadAudio(ctx, fileBytes)
}

func (s *transcribeService) TranscribeAudio(ctx context.Context, uploadedFileUrl string) (string, error) {
	return s.transcriptionClient.TranscribeAudio(ctx, uploadedFileUrl)
}