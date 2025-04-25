package app

import (
	"github.com/ezydark/ezHead/libs"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/rs/zerolog/log"
)

var rodBrowser *rod.Browser

// InitRodBrowser initializes the RodBrowser instance.
func InitRodBrowser() *rod.Browser {
	if rodBrowser != nil {
		log.Fatal().Msg("RodBrowser already initialized")
	}

	u := launcher.New().
		Headless(true).
		UserDataDir(libs.GetAppFolder().String()).
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

	rodBrowser = rod.New().ControlURL(u).NoDefaultDevice().MustConnect()
	return rodBrowser
}

func GetRodBrowser() *rod.Browser {
	if rodBrowser == nil {
		log.Fatal().Msg("RodBrowser not initialized")
	}
	return rodBrowser
}
