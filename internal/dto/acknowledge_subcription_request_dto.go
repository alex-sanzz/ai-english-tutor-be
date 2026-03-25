package dto


type AcknowledgeSubcriptionRequestDto struct {
	SubscriptionId string `json:"subscription_id"`
	PurchaseToken string `json:"purchase_token"`
}