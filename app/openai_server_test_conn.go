package app

import (
	openai "github.com/ezydark/ezHead/libs/openai/client"
	"github.com/rs/zerolog/log"
	go_openai "github.com/sashabaranov/go-openai"
)

func OpenAIServerTestConnection() go_openai.ChatCompletionResponse {
	client := openai.InitClient()
	response, err := client.TestConnection()

	if err != nil {
		log.Fatal().Msgf("OpenAI server test connection failed:\n%v", err)
	}

	return response
}
