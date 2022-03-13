package config

import (
	"mangathrV2/internal/sources/connections/mangadex"
)

type Config struct {
	Database struct {
		Driver   string
		Uri      string
	}
	Metadata struct {
		Agent             string
		Location          string
	}
	Downloader struct {
		SimultaneousPages int `yaml:"simultaneousPages"`
		PageRetries       int `yaml:"pageRetries"`
		Delay struct {
			Page          int
			Chapter       int
			UpdateChapter int `yaml:"updateChapter"`
		}
		Output struct {
			Path          string
			UpdatePath    string `yaml:"updatePath"`
			Zip           bool
		}
	}
	Sources struct {
		Mangadex mangadex.Config
	}

}
