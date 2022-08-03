package defaults

import (
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/ui"
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

	configDir := filepath.Join(homeDir, ".config", "mangathrv2")
	return configDir
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config")
}

/*
	Config defaults
*/

func DatabaseDriver() string {
	return "sqlite"
}

func DatabaseUri() string {
	return filepath.Join(ConfigDir(), "db.sqlite")
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

func DatabaseUriDocker() string {
	return filepath.Join(ConfigDirDocker(), "db.sqlite")
}

func DataPathDocker() string {
	return "/data"
}
