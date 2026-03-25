package androidpublisher

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/config"
	"context"
	"fmt"

	publisher "google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

type AndroidPublisherClient interface {
    VerifySubscription(ctx context.Context, subscriptionId, purchaseToken string) (*publisher.SubscriptionPurchase, error)
    AcknowledgeSubscription(ctx context.Context, subscriptionId, purchaseToken string) error
}

type androidPublisherClient struct {
	service *publisher.Service
	cfg config.AppConfig
}

func New(ctx context.Context, credentialPath string) (AndroidPublisherClient, error) {
	service, err := publisher.NewService(ctx, option.WithCredentialsFile(credentialPath))

	if err != nil {
		return nil, apperr.Internal(fmt.Errorf("android publisher client new error : %w", err))
	}

	return androidPublisherClient{
		service: service,
	}, nil
}

func (c androidPublisherClient) VerifySubscription(ctx context.Context, subscriptionId, purchaseToken string) (*publisher.SubscriptionPurchase, error){
	result, err := c.service.Purchases.Subscriptions.Get(c.cfg.PackageName, subscriptionId, purchaseToken).Do()

	if err != nil {
		return nil, fmt.Errorf("android publisher client verify subscription error: %w", err)
	}

	return result, nil 
}

func (c androidPublisherClient) AcknowledgeSubscription(ctx context.Context, subscriptionId, purchaseToken string) error{
	err := c.service.Purchases.Subscriptions.Acknowledge(c.cfg.PackageName, subscriptionId, purchaseToken, &publisher.SubscriptionPurchasesAcknowledgeRequest{}).Do()

	if err != nil {
		return fmt.Errorf("android publisher client verify subscription error: %w", err)
	}

	return nil
}