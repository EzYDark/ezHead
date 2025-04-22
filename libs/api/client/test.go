package client

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

func TestServer() error {
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
		return fmt.Errorf("error creating chat completion:\n%v", err)
	}

	log.Debug().Msgf("Test Response: %s\n", resp.Choices[0].Message.Content)

	return nil
}
