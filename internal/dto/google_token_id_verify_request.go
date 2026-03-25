package dto

type GoogleTokenIdVerifyRequest struct {
	GoogleTokenId string `json:"google_token_id"`
	Nonce string `json:"nonce"`
}