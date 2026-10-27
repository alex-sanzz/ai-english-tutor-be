package service

import (
	"ai-tutor-backend/internal/apperr"
	androidpublisher "ai-tutor-backend/internal/infrastructure/google/android-publisher"
	"ai-tutor-backend/internal/models"
	"ai-tutor-backend/internal/repository"
	"context"
	"fmt"
	"time"
)

type paymentService struct {
	userRepository repository.UserRepository
	androidPublisherClient androidpublisher.AndroidPublisherClient
	subcriptionRepository repository.SubcriptionRepository
}

func NewPaymentService(androidPublisherClient androidpublisher.AndroidPublisherClient, userRepository repository.UserRepository, subcriptionRepository repository.SubcriptionRepository) paymentService{
	return paymentService{
		androidPublisherClient: androidPublisherClient,
		userRepository: userRepository,
		subcriptionRepository: subcriptionRepository,
	}
}

func (p paymentService) AcknowledgeSubscription(ctx context.Context, userId, subscriptionId, purchaseToken string) error{
	
	// _, err := p.userRepository.FindById(ctx, userId)

	// if err != nil {
	// 	return false, nil
	// }

	s, err := p.subcriptionRepository.FindActiveSubscription(ctx, userId)

	if err != nil {
		return err
	}

	if s != nil {
		return apperr.BadRequest("400", "previous subcriptions is still active", fmt.Errorf("payment service error: acknowledge subcription is still active"))
	}

	purchase, err := p.androidPublisherClient.VerifySubscription(ctx, subscriptionId, purchaseToken)

	if err != nil {
		return apperr.Internal(fmt.Errorf("payment service impl verify subscription error: %w", err)) 
	}

	item := purchase.LineItems[0]

	expiryAt, err := time.Parse(time.RFC3339Nano, item.ExpiryTime)

	if err != nil {
		return apperr.Internal(fmt.Errorf("payment service impl time parse error: %w", err))
	}

	err = p.subcriptionRepository.Create(ctx, models.Subcription{
		UserId: userId,
		PurchaseToken: purchaseToken,
		PurchaseTime: time.Now(),
		ProductId: item.ProductId,
		BasePlanId: item.OfferDetails.BasePlanId,
		ExpiryAt: expiryAt,
		AutoRenewing: item.AutoRenewingPlan != nil && purchase.LineItems[0].AutoRenewingPlan.AutoRenewEnabled,
	})

	if err != nil {
		return apperr.Internal(fmt.Errorf("payment service impl verify subscription error: %w", err)) 
	}
	
	p.androidPublisherClient.AcknowledgeSubscription(ctx, subscriptionId, purchaseToken)

	return nil 
}

