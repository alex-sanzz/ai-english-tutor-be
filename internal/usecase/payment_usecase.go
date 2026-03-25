package usecase

import (
	"ai-tutor-backend/internal/service"
	"context"
)

type PaymentUseCase struct {
	paymentService service.PaymentService
}

func NewPaymentUseCase(paymentService service.PaymentService) *PaymentUseCase{
	return &PaymentUseCase{
		paymentService: paymentService,
	}
}

func (u *PaymentUseCase) AcknowledgeSubcription(ctx context.Context, userId, subcriptionId, purchaseToken string) error{
	return u.paymentService.AcknowledgeSubscription(ctx, userId, subcriptionId, purchaseToken)
}

