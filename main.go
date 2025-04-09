package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ezydark/ezHead/libs/perplexity"
	"github.com/ezydark/ezforce/libs/logger"
	"github.com/fatih/color"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize logger
	err := logger.Init()
	if err != nil {
		fatal_tag := color.New(color.FgRed, color.Bold).Sprintf("[FATAL]")
		fmt.Println(fatal_tag, "Could not initialize custom logger:\n", err)
		return
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		// Fallback for older Windows versions
		userProfile := os.Getenv("USERPROFILE")
		localAppData = filepath.Join(userProfile, "AppData", "Local")
	}
	userDataDir := filepath.Join(localAppData, "ezHead")

	u := launcher.New().
		Headless(false).
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

	// Execute custom request
	headerJSON, err := perplexity.DefaultHeaders().IntoJSON()
	if err != nil {
		log.Fatal().Msgf("Could not convert headers to JSON:\n%v", err)
	}
	bodyJSON, err := perplexity.DefaultBody().IntoJSON()
	if err != nil {
		log.Fatal().Msgf("Could not convert body to JSON:\n%v", err)
	}
	jsScript := perplexity.DefaultJsScript(headerJSON, bodyJSON)
	page.MustElement(`textarea[placeholder="Ask anything..."]`)
	result := page.MustEval(jsScript)
	resultJSON, err := result.MarshalJSON()
	if err != nil {
		log.Fatal().Msgf("Could not convert result to JSON:\n%v", err)
	}
	var perplexityResponse perplexity.Response
	if err := json.Unmarshal(resultJSON, &perplexityResponse); err != nil {
		prettyJSON := result.JSON("", "  ")
		log.Fatal().Msgf("Error parsing result:\n%v\nRaw response:\n%s", err, prettyJSON)
	}
	finalAnswer, err := perplexityResponse.FinalMessage.GetFinalAnswer()
	if err != nil {
		log.Error().Msgf("Error getting final answer: %v", err)
	} else if finalAnswer != "" {
		log.Info().Msgf("Final answer: %s", finalAnswer)
	}

	time.Sleep(time.Hour)
}
