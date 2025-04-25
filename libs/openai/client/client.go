package openai

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	go_openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client    *go_openai.Client
	authToken string
	baseUrl   string
}

// Default OpenAI-compatible client configuration
var (
	def_authToken = "ez-123456789"
	def_baseUrl   = "http://localhost:8080/v1"
)

// InitClient initializes an OpenAIClient with default settings
func InitClient() *OpenAIClient {
	config := go_openai.DefaultConfig(def_authToken)
	config.BaseURL = def_baseUrl

	return &OpenAIClient{
		client:    go_openai.NewClientWithConfig(config),
		authToken: def_authToken,
		baseUrl:   def_baseUrl,
	}
}

// TestConnection tests the connection to the OpenAI-compatible server
// (NOTICE: Only simple non-streaming connection)
func (oai_client *OpenAIClient) TestConnection() (go_openai.ChatCompletionResponse, error) {
	resp, err := oai_client.client.CreateChatCompletion(
		context.Background(),
		go_openai.ChatCompletionRequest{
			Model: "auto",
			Messages: []go_openai.ChatCompletionMessage{
				{
					Role:    go_openai.ChatMessageRoleUser,
					Content: "Hello from local OpenAI-compatible test client!",
				},
			},
		},
	)

	if err != nil {
		return go_openai.ChatCompletionResponse{}, fmt.Errorf("error creating chat completion:\n%v", err)
	}

	log.Debug().Msgf("Response from OpenAI-compatible client:\n%s", resp.Choices[0].Message.Content)
	return resp, nil
}
