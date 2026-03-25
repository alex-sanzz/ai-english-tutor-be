package auth

import (
	"context"

	"cloud.google.com/go/auth/credentials/idtoken"
)

type GoogleAuthClient interface {
	ParseAndVerifyGoogleTokenId(context context.Context, googleTokenId string, nonce string) (*idtoken.Payload, error)
}