package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ezydark/ezHead/app"
	"github.com/ezydark/ezHead/libs"
	"github.com/ezydark/ezHead/libs/perplexity"
	"github.com/rs/zerolog/log"
)

func main() {
	app.InitLogger()
	_ = libs.InitAppFolder()
	_ = app.InitRodBrowser()
	_ = app.InitOpenAIServer()
	_ = app.OpenAIServerTestConnection()

	// perplexity.PerplexPage = rod.New().ControlURL(u).NoDefaultDevice().MustConnect().MustPage("https://www.perplexity.ai/")
	// perplexity.PerplexPage.MustWindowMaximize()
	// perplexity.PerplexPage.Browser().SlowMotion(1 * time.Second)
	// perplexity.PerplexPage.Browser().Trace(true)

	// perplexity.PerplexPage.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
	// 	UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36 Edg/135.0.0.0",
	// })
	// perplexity.PerplexPage.MustSetViewport(1920, 1080, 1.0, false)

	// _ = perplexity.PerplexPage.MustExpose("goLogInfo", func(g gson.JSON) (any, error) {
	// 	log.Info().Msgf("[JS] %v", g.Str())
	// 	return nil, nil
	// })

	// _ = perplexity.PerplexPage.MustExpose("goProcessStreamChunk", func(chunk gson.JSON) (any, error) {
	// 	_, err := response.ProcessStreamChunk(chunk)
	// 	if err != nil {
	// 		log.Fatal().Msgf("Error processing stream chunk:\n%v", err)
	// 	}

	// 	return nil, nil
	// })

	// _ = perplexity.PerplexPage.MustExpose("goLogError", func(g gson.JSON) (any, error) {
	// 	log.Error().Msgf("[JS] %v", g.Str())
	// 	return nil, nil
	// })

	// _ = perplexity.PerplexPage.MustExpose("goLogFatal", func(g gson.JSON) (any, error) {
	// 	log.Fatal().Msgf("[JS] %v", g.Str())
	// 	return nil, nil
	// })

	// perplex, err := perplexity.Init()
	// if err != nil {
	// 	log.Fatal().Msgf("Could not initialize Perplexity struct:\n%v", err)
	// }

	for {
		fmt.Println("\n\nEnter new query to send to Perplexity >> ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			query := scanner.Text()

			err := perplex.SendRequest(perplexity.PerplexPage, query)
			if err != nil {
				log.Fatal().Msgf("Could not send request to Perplexity:\n%v", err)
			}

			// finalAnswer, err := res.FinalMessage.GetFinalAnswer()
			// if err != nil {
			// 	log.Fatal().Msgf("Could not get final answer from Perplexity:\n%v", err)
			// }
			// log.Info().Msgf("Final Answer:\n%s", finalAnswer)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal().Err(err).Msg("Error reading input")
		}

		// time.Sleep(time.Hour)
	}
}
