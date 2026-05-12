package openai

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/config"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/openai/openai-go/v3"

)

type OpenAiClient struct {
	client openai.Client
	config *config.ConfigApp
	logger log.Logger
}

func NewOpenAiClient(c openai.Client, config *config.ConfigApp, logger log.Logger) *OpenAiClient {
	return &OpenAiClient{
		client: c,
		config: config,
		logger: logger,
	}
}

func (c *OpenAiClient) AskQuestion(ctx context.Context, systemPrompt, question, answer string) (string, error){

	openAiMessages := make([]openai.ChatCompletionMessageParamUnion, 0, 3)

	if systemPrompt != "" {
		openAiMessages = append(openAiMessages, openai.ChatCompletionMessageParamUnion{
			// if the chat comes from user then you should use ofUser
			// But if it comes from an AI, then you should use OfAssistant
			OfSystem: &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: openai.String(systemPrompt),
				},
			},
		})
	}

	if question != "" {

		openAiMessages = append(openAiMessages, openai.ChatCompletionMessageParamUnion{
			// if the chat comes from user then you should use ofUser
			// But if it comes from an AI, then you should use OfAssistant
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				Content: openai.ChatCompletionAssistantMessageParamContentUnion{
					OfString: openai.String(question),
				},
			},
		})
	}

	openAiMessages = append(openAiMessages, openai.ChatCompletionMessageParamUnion{
		// if the chat comes from user then you should use ofUser
		// But if it comes from an AI, then you should use OfAssistant
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: openai.String(answer),
			},
		},
	})
	
	

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
        Model: c.config.OpenAi.Model,
        Messages: openAiMessages,
    })

	if err != nil {
		return "", apperr.Internal(err)
	}

	if len(resp.Choices) == 0 {
        return "", apperr.Internal(errors.New("no response from OpenAI"))
    }

	return resp.Choices[0].Message.Content, nil

	
}

func (c *OpenAiClient) ChatStream(ctx context.Context, systemPrompt string, message []*models.Chat, onChunk func(string) error, onFinish func(string) error) error{
	
	
	openAiMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(message))

	openAiMessages = append(openAiMessages, openai.ChatCompletionMessageParamUnion{
		// if the chat comes from user then you should use ofUser
		// But if it comes from an AI, then you should use OfAssistant
		OfSystem: &openai.ChatCompletionSystemMessageParam{
			Content: openai.ChatCompletionSystemMessageParamContentUnion{
				OfString: openai.String(systemPrompt),
			},
		},
	})
	
	for _, msg := range message {
		switch msg.Role {
		case "user":
			openAiMessages = append(openAiMessages, openai.ChatCompletionMessageParamUnion{
				// if the chat comes from user then you should use ofUser
				// But if it comes from an AI, then you should use OfAssistant
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(msg.Message),
					},
				},
			})
		case "assistant":
			openAiMessages = append(openAiMessages, openai.ChatCompletionMessageParamUnion{
				// if the chat comes from user then you should use ofUser
				// But if it comes from an AI, then you should use OfAssistant
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String(msg.Message),
					},
				},
			})
		}
	}

	stream := c.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model: c.config.OpenAi.Model,
		// This messages is actually a chat histories
		// So all of previous messages should be inserted here
		// Union means the chat can have different kind
		// For example: it can be from human, ai, etc
		Messages: openAiMessages,
		
	})

	if stream == nil {
		c.logger.Error("stream is nil")
		return errors.New("openai client chat stream: streaming returns nil")
	}

	fullSentences := ""
	defer stream.Close()
	
	for stream.Next() {
		event := stream.Current()

		for _, delta := range event.Choices{
			if err := onChunk(delta.Delta.Content); err != nil {
				c.logger.Debug("openai client chat stream onChunk error")
				return fmt.Errorf("openai client chat stream: %w", err)
			}
			fullSentences += delta.Delta.Content
		}
	}

	if err := stream.Err(); err != nil {
		if !errors.Is(err, context.Canceled){
			c.logger.Debug("openai client chat stream error (stream.Err() not nil)")
			return fmt.Errorf("openai client chat stream: %w", err)
		}
		
	}
	

	if err := onFinish(fullSentences); err != nil {
		c.logger.Debug("openai client chat stream error (onFinish() not nil)")
		return fmt.Errorf("openai client chat stream: %w", err)
	}
	

	return nil

	
}