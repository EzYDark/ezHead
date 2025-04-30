package libs

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/rs/zerolog/log"
	"github.com/ysmood/gson"
)

// Single instance of RodBrowser
var rodBrowser *rod.Browser

// InitRodBrowser initializes the RodBrowser instance.
func InitRodBrowser() *rod.Browser {
	if rodBrowser != nil {
		log.Fatal().Msg("RodBrowser already initialized")
	}

	u := launcher.New().
		Headless(true).
		UserDataDir(GetAppFolder().String()).
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

	// For debugging purposes
	// rodBrowser.SlowMotion(1 * time.Second)
	// rodBrowser.Trace(true)

	return rodBrowser
}

func GetRodBrowser() *rod.Browser {
	if rodBrowser == nil {
		log.Fatal().Msg("RodBrowser not initialized")
	}
	return rodBrowser
}

func SetPageSettings(page *rod.Page) *rod.Page {
	page.MustWindowMaximize()
	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36 Edg/135.0.0.0",
	})
	page.MustSetViewport(1920, 1080, 1.0, false)

	return page
}

func ExposeGoLogger(page *rod.Page) *rod.Page {
	_ = page.MustExpose("goLogInfo", func(g gson.JSON) (any, error) {
		log.Info().Msgf("[JS] %v", g.Str())
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
	return page
}
