package downloader

import (
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/ui"
	"os"
	"path/filepath"
)

type Config struct {
	DryRun            bool `yaml:"dryRun"`
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

func (c *Config) Default() {
	c.DryRun = false
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
}

func getCWD() string {
	path, err := os.Getwd()
	if err != nil {
		logging.Errorln(err)
		ui.Fatal("Failed to find current working directory.")
	}
	return filepath.Join(path, "mangathrv2")
}
