package handler

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/dto"
	"ai-tutor-backend/internal/usecase"
	"fmt"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUseCase usecase.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase usecase.PaymentUseCase) *PaymentHandler{
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

func (h *PaymentHandler) AcknowledgeSubcription(c *gin.Context){
	var p dto.AcknowledgeSubcriptionRequestDto

	if err := c.BindJSON(&p); err != nil {
		writeError(c, apperr.BadRequest("400", "wrong JSON body format", fmt.Errorf("payment handler error: cannot map request body, error: %w", err)))
		return
	}

	userId := c.GetString("userId")

	if userId == "" {
		writeError(c, apperr.Unauthorized("401", "Unauthorized", fmt.Errorf("payment handler error: cannot get user id from jwt token")))
		return
	}

	err := h.paymentUseCase.AcknowledgeSubcription(c.Request.Context(), userId, p.SubscriptionId, p.PurchaseToken)

	if err != nil {
		writeError(c, apperr.Internal(err))
		return
	}

	c.JSON(204, nil)
}