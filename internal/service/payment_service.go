package service

import "context"

type PaymentService interface {
	AcknowledgeSubscription(ctx context.Context, userId, subscriptionId, purchaseToken string) error
	
}