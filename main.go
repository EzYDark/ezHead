package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ezydark/ezHead/libs/api/client"
	"github.com/ezydark/ezHead/libs/api/server"
	"github.com/ezydark/ezHead/libs/perplexity"
	"github.com/ezydark/ezHead/libs/perplexity/response"
	"github.com/ezydark/ezforce/libs/logger"
	"github.com/fatih/color"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/rs/zerolog/log"
	"github.com/ysmood/gson"
)

func main() {
	// Initialize logger
	err := logger.Init()
	if err != nil {
		fatal_tag := color.New(color.FgRed, color.Bold).Sprintf("[FATAL]")
		fmt.Println(fatal_tag, "Could not initialize custom logger:\n", err)
		return
	}

	http.HandleFunc("/v1/chat/completions", server.AuthMiddleware(server.HandleChatCompletions))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Error().Msgf("Failed to write health response:\n%v", err)
		}
	})

	log.Info().Msg("OpenAI-compatible server starting on port 8080...")
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- http.ListenAndServe(":8080", nil)
	}()
	go func() {
		err := <-serverErrors
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Server error while starting:\n%v", err)
		} else if err == http.ErrServerClosed {
			log.Info().Msg("Server closed")
		}
	}()

	err = client.TestServer()
	if err != nil {
		log.Fatal().Msgf("Failed to test API server:\n%v", err)
	}

	// time.Sleep(time.Hour)

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		// Fallback for older Windows versions
		userProfile := os.Getenv("USERPROFILE")
		localAppData = filepath.Join(userProfile, "AppData", "Local")
	}
	userDataDir := filepath.Join(localAppData, "ezHead")

	u := launcher.New().
		Headless(true).
		UserDataDir(userDataDir).
		ProfileDir("Default").
		NoSandbox(true).
		Set("disable-extensions", "false").
		Set("enable-automation", "false").
		Set("disable-features", "IsolateOrigins,site-per-process").
		Set("disable-web-security", "true").
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-sync", "true").
		// Set("load-extension", "path\\to\\extension").
		MustLaunch()

	page := rod.New().ControlURL(u).NoDefaultDevice().MustConnect().MustPage("https://www.perplexity.ai/")
	page.MustWindowMaximize()
	page.Browser().SlowMotion(1 * time.Second)
	page.Browser().Trace(true)

	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36 Edg/135.0.0.0",
	})
	page.MustSetViewport(1920, 1080, 1.0, false)

	_ = page.MustExpose("goLogInfo", func(g gson.JSON) (any, error) {
		log.Info().Msgf("[JS] %v", g.Str())
		return nil, nil
	})

	_ = page.MustExpose("goProcessStreamChunk", func(chunk gson.JSON) (any, error) {
		_, err := response.ProcessStreamChunk(chunk)
		if err != nil {
			log.Fatal().Msgf("Error processing stream chunk:\n%v", err)
		}

		return nil, nil
	})

	_ = page.MustExpose("goLogError", func(g gson.JSON) (any, error) {
		log.Error().Msgf("[JS] %v", g.Str())
		return nil, nil
	})

	_ = page.MustExpose("goLogFatal", func(g gson.JSON) (any, error) {
		log.Fatal().Msgf("[JS] %v", g.Str())
		return nil, nil
	})

	perplex, err := perplexity.Init()
	if err != nil {
		log.Fatal().Msgf("Could not initialize Perplexity struct:\n%v", err)
	}

	for {
		fmt.Println("\n\nEnter new query to send to Perplexity >> ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			query := scanner.Text()

			err := perplex.SendRequest(page, query)
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
