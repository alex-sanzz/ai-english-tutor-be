package jwt

import (
	"ai-tutor-backend/internal/config"
	"ai-tutor-backend/internal/log"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClient interface {
	GenerateAccessToken(userId string) (string, error)
	GenerateRefreshToken() (string, error)
	ParseAndValidate(token string) (*jwt.RegisteredClaims, error)
}

type rsaJwtClient struct {
	cfg config.JwtConfig
	privateKey *rsa.PrivateKey
	publicKey *rsa.PublicKey
	logger log.Logger
}

func NewRsaJwtClient(cfg config.JwtConfig, logger log.Logger) (*rsaJwtClient, error){
	privateKeyBytes, err := os.ReadFile(cfg.PrivateKeyPath)

	if err != nil {
		return nil, fmt.Errorf("jwt client new rsa jwt client error: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	if err != nil {
		return nil, fmt.Errorf("jwt client new rsa jwt client error: %w", err)
	}

	publicKeyBytes, err := os.ReadFile(cfg.PublicKeyPath)

	if err != nil {
		return nil, fmt.Errorf("jwt client new rsa jwt client error: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)

	if err != nil {
		return nil, fmt.Errorf("jwt client new rsa jwt client error: %w", err)
	}

	return &rsaJwtClient{
		cfg: cfg,
		privateKey: privateKey,
		publicKey: publicKey,
		logger: logger,
	}, nil 
}

func (c *rsaJwtClient) GenerateAccessToken(userId string) (string, error){
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject: userId,
		Issuer: c.cfg.Issuer,
		Audience: []string{c.cfg.Audience},
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(c.cfg.AccessTokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(c.privateKey)
}

func (c *rsaJwtClient) ParseAndValidate(token string) (*jwt.RegisteredClaims, error){
	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error){
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("jwt client parse and validate error: unexpected signing method")
		}
		return c.publicKey, nil 
	})

	if err != nil {
		return nil, fmt.Errorf("jwt client parse and validate error: %w", err) 
	}

	claims, ok := parsed.Claims.(*jwt.RegisteredClaims)

	if !ok || !parsed.Valid {
		return nil, fmt.Errorf("jwt client parse and validate error: invalid token") 
	}

	if claims.Issuer != c.cfg.Issuer {
		return nil, fmt.Errorf("jwt client parse and validate error: invalid issuer")
	}

	return claims, nil 
}

func (c *rsaJwtClient) GenerateRefreshToken() (string, error){

	randomBytes := make([]byte, 32)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("jwt client generate refresh token error: %w", err)
	}

	hashedToken := hashToken(string(randomBytes))

	return hashedToken, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))

	return base64.URLEncoding.EncodeToString(hash[:])
}

