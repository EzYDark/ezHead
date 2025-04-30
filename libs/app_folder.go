package libs

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type AppFolderPath string

// Single instance of AppFolderPath
var appFolder *AppFolderPath

// Initialize the app's default folder
func InitAppFolder() *AppFolderPath {
	if appFolder == nil {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			// Fallback for older Windows versions
			userProfile := os.Getenv("USERPROFILE")
			localAppData = filepath.Join(userProfile, "AppData", "Local")
		}
		path := AppFolderPath(filepath.Join(localAppData, "ezHead"))
		appFolder = &path
	}
	return appFolder
}

func GetAppFolder() *AppFolderPath {
	if appFolder == nil {
		log.Fatal().Msg("AppFolder is not initialized")
	}

	return appFolder
}

// Set the app's folder to a own folder path
func (af *AppFolderPath) Set(newFolderPath string) *AppFolderPath {
	*af = AppFolderPath(newFolderPath)
	return af
}

func (af *AppFolderPath) String() string {
	if af == nil {
		log.Fatal().Msg("AppFolder is not initialized")
		return "" // Satisfy static analyzer
	}

	return string(*af)
}
