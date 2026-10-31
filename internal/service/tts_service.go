package service

import (
	"context"
	"io"
)

type TTSService interface {
	Generate(ctx context.Context, text string) (io.ReadCloser, error)
}