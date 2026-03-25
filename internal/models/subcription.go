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
	ExpiryAt time.Time `json:"expiry_at"`
	AutoRenewing bool `json:"auto_renewing"`
	CreatedAt time.Time `json:"created_at"`
}