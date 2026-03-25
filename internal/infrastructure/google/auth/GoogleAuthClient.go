package auth

import (
	"ai-tutor-backend/internal/config"
	"ai-tutor-backend/internal/log"
	"context"
	"fmt"

	"cloud.google.com/go/auth/credentials/idtoken"
	"go.uber.org/zap"
)

type googleAuthClient struct {
	config *config.GoogleAuthConfig
	logger log.Logger
}

func NewGoogleAuthCLient(config *config.GoogleAuthConfig, logger log.Logger) *googleAuthClient {
	return &googleAuthClient{
		config: config,
		logger: logger,
	}
}

func (g *googleAuthClient) ParseAndVerifyGoogleTokenId(context context.Context, googleTokenId string, nonce string) (*idtoken.Payload, error) {
	// This automatically gets google public key to verify the google token id
	payload, err := idtoken.Validate(context, googleTokenId, g.config.ClientId)

	//  if you check jwt field, there is aud, azp, sub, email, etc
	// aud is equal to google client id
	// sub is user id, user email might be changed, but even though it's changed, sub remains unchanged 
	g.logger.Debug("google token id payload", zap.String("google token id", googleTokenId))

	if err != nil {
		return nil, fmt.Errorf("google auth client verify google token id error : %w", err) 
	}

	nonceFromToken := payload.Claims["nonce"].(string)

	if nonce != nonceFromToken {
		return nil, fmt.Errorf("google auth client verify google token id error : nonce is invalid")
	}

	return payload, nil


}	