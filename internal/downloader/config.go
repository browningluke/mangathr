package downloader

import (
	"github.com/browningluke/mangathr/internal/config/defaults"
	"github.com/browningluke/mangathr/internal/logging"
	"github.com/browningluke/mangathr/internal/ui"
	"os"
	"path/filepath"
)

var config Config

type Config struct {
	DryRun            bool `yaml:"dryRun"`
	CleanupOnError    bool `yaml:"cleanupOnError"`
	SimultaneousPages int  `yaml:"simultaneousPages"`
	PageRetries       int  `yaml:"pageRetries"`
	Delay             struct {
		Page          string
		Chapter       string
		UpdateChapter string `yaml:"updateChapter"`
	}
	Output struct {
		Path             string
		UpdatePath       string `yaml:"updatePath"`
		Zip              bool
		FilenameTemplate string `yaml:"filenameTemplate"`
	}
	Metadata struct {
		Agent    string
		Location string
	}
}

func SetConfig(cfg Config) {
	config = cfg
}

func DryRun() bool {
	return config.DryRun
}

func (c *Config) Default(inContainer bool) {
	c.DryRun = false
	c.CleanupOnError = true
	c.SimultaneousPages = 2
	c.PageRetries = 5

	c.Delay.Page = "50ms"
	c.Delay.Chapter = "100ms"
	c.Delay.UpdateChapter = "250ms"

	c.Output.Path = getCWD()       // Use CWD
	c.Output.UpdatePath = getCWD() // Use CWD
	c.Output.Zip = true
	c.Output.FilenameTemplate = "{num:3} - Chapter {num}{title: - <.>}{groups: [<.>]}"

	c.Metadata.Agent = "comicinfo"
	c.Metadata.Location = "internal"

	// Overwrite defaults if we are in a container
	if inContainer {
		c.Output.Path = defaults.DataPathDocker()
		c.Output.UpdatePath = defaults.DataPathDocker()
	}
}

func getCWD() string {
	path, err := os.Getwd()
	if err != nil {
		logging.Errorln(err)
		ui.Fatal("Failed to find current working directory.")
	}
	return filepath.Join(path, "mangathr")
}
