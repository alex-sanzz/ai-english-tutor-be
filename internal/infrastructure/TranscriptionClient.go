package infrastructure

import "context"

type TranscriptionClient interface {
    UploadAudio(ctx context.Context, fileBytes []byte) (string, error)
    TranscribeAudio(ctx context.Context, uploadedFileUrl string) (string, error)
}