package main

import (
	"time"

	"github.com/ezydark/ezHead/app"
	"github.com/ezydark/ezHead/libs"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting ezHead...")

	app.InitLogger()
	_ = libs.InitAppFolder()

	_ = libs.InitRodBrowser()

	oai_server := app.InitOpenAIServer()
	oai_server.Start()
	_ = app.OpenAIServerTestConnection()

	log.Info().Msg("Sleeping main thread for 999 hours...")
	time.Sleep(999 * time.Hour)
}
