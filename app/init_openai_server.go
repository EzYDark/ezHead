package app

import (
	openai "github.com/ezydark/ezHead/libs/openai/server"
	"github.com/rs/zerolog/log"
)

func InitOpenAIServer() *openai.OpenAIServer {
	server, err := openai.InitServer()
	if err != nil {
		log.Fatal().Msgf("Failed to initialize OpenAI server:\n%v", err)
	}
	return server
}
