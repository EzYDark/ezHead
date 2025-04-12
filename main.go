package main

import (
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
	page.MustSetViewport(1920, 1080, 1.0, false)

	perplex, err := perplexity.Init()
	if err != nil {
		log.Fatal().Msgf("Could not initialize Perplexity struct:\n%v", err)
	}
	res, err := perplex.SendRequest(page, "Is this working, lul? <3")
	if err != nil {
		log.Fatal().Msgf("Could not send request to Perplexity:\n%v", err)
	}

	finalAnswer, err := res.FinalMessage.GetFinalAnswer()
	if err != nil {
		log.Fatal().Msgf("Could not get final answer from Perplexity:\n%v", err)
	}
	log.Info().Msgf("Final Answer:\n%s", finalAnswer)

	// time.Sleep(time.Hour)
}
