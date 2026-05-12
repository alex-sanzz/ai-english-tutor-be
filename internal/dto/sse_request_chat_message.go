package dto


type SseChatMessageRequest struct {
	Message     string `json:"message"`
	UserId 		string   `json:"user_id"`
	Type        string `json:"type"`
	CulturalContext string `json:"cultural_context"`
	AlternateVersion string `json:"alternate_version"`
	ParaphraseVersion string `json:"paraphrase_version"`
	
}
