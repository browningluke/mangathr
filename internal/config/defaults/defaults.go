package defaults

import (
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/ui"
	"os"
	"path/filepath"
)

/*
	Config directory defaults
*/

func ConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logging.Errorln(err)
		ui.Fatal("Failed to find config directory.")
	}

	configDir := filepath.Join(homeDir, ".config", "mangathr")
	return configDir
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config")
}

/*
	Container defaults
*/

func ConfigDirDocker() string {
	configDir := "/config"
	return configDir
}

func ConfigPathDocker() string {
	return filepath.Join(ConfigDirDocker(), "config")
}

func DataPathDocker() string {
	return "/data"
}
