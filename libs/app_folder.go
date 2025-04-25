package libs

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type AppFolderPath string

var appFolder *AppFolderPath

// Initialize the app's default folder
func InitAppFolder() *AppFolderPath {
	if *appFolder == "" {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			// Fallback for older Windows versions
			userProfile := os.Getenv("USERPROFILE")
			localAppData = filepath.Join(userProfile, "AppData", "Local")
		}
		*appFolder = AppFolderPath(filepath.Join(localAppData, "ezHead"))
	}
	return appFolder
}

func GetAppFolder() *AppFolderPath {
	if *appFolder == "" {
		log.Fatal().Msg("AppFolder is not initialized")
	}

	return appFolder
}

// Set the app's folder to a own folder path
func (af *AppFolderPath) Set(newFolderPath string) {
	*af = AppFolderPath(newFolderPath)
}

func (af *AppFolderPath) String() string {
	if *af == "" {
		log.Fatal().Msg("AppFolder is not initialized")
	}

	return string(*af)
}
