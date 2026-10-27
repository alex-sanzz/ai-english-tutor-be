package models

import (
	"time"

	"github.com/google/uuid"
)

type Subcription struct {
	ID uuid.UUID `json:"id"`
	UserId string `json:"user_id"`
	PurchaseToken string `json:"purchase_token"`
	PurchaseTime time.Time `json:"purchase_time"`
	// Example:
	// productId -> premium
	// baseplanid -> monthly
	// This shows which product/subscription the user paid
	ProductId string `json:"product_id"`
	// This shows based on the product/subscription, how is it billed
	BasePlanId string `json:"base_plan_id"`
	ExpiryAt time.Time `json:"expiry_at"`
	AutoRenewing bool `json:"auto_renewing"`
	CreatedAt time.Time `json:"created_at"`
}