package client

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

func TestServer() {
	config := openai.DefaultConfig("your-test-token")
	config.BaseURL = "http://localhost:8080/v1" // Point to your local server

	client := openai.NewClientWithConfig(config)

	// Test with an actual request
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "your-local-model",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello from local test!",
				},
			},
		},
	)

	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	log.Info().Msgf("Test Response: %s\n", resp.Choices[0].Message.Content)
}
