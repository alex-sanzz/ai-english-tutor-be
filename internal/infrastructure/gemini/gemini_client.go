package gemini

import (
	"context"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient(client *genai.Client) *GeminiClient{
	return &GeminiClient{
		client: client,
	}
}

func (c *GeminiClient) Chat(ctx context.Context, message string)(string, error){
	result, err := c.client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(message), nil)
	
	if err != nil {
		return "", err 
	}

	return result.Text(), nil
}

func (c *GeminiClient) ChatStream(ctx context.Context, message string, onChunk func(string) error, onFinish func() error) error{
	stream := c.client.Models.GenerateContentStream(ctx, "gemini-2.5-flash", genai.Text(message), nil)
	
	for chunk, err := range stream {
		if err != nil {
			return err
		}
		part := chunk.Candidates[0].Content.Parts[0]
      	onChunk(part.Text)
	}
	onFinish()

	return nil
}