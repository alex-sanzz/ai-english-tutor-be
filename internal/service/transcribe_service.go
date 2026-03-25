package service

import "context"

type TranscribeService interface {
	UploadAudioFile(ctx context.Context, fileBytes []byte) (string, error)
	TranscribeAudio(ctx context.Context, uploadedFileUrl string) (string, error)
}