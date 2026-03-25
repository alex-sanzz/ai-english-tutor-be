package usecase

import (
	"ai-tutor-backend/internal/service"
	"context"
)

type TranscribeUseCase struct {
	transcriptionService service.TranscribeService
}

func NewTranscribeUseCase(transcription service.TranscribeService) *TranscribeUseCase {
	return &TranscribeUseCase{
		transcriptionService: transcription,
	}
}

func (u *TranscribeUseCase) UploadAudioFile(ctx context.Context, fileBytes []byte) (string, error) {
	return u.transcriptionService.UploadAudioFile(ctx, fileBytes)
}

func (u *TranscribeUseCase) TranscribeAudio(ctx context.Context, uploadedFileUrl string) (string, error)  {
	return u.transcriptionService.TranscribeAudio(ctx, uploadedFileUrl)
}