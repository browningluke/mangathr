package mangaplus

import (
	"errors"
	"github.com/browningluke/mangathr/v2/internal/utils"
)

var config Config

type Config struct {
	FilenameTemplate string `yaml:"filenameTemplate"`
	Language         int    `yaml:"language"`
	ImageQuality     string `yaml:"imageQuality"`
	Split            string `yaml:"split"`
}

func SetConfig(cfg Config) {
	config = cfg
}

func (c *Config) Default() {
	c.FilenameTemplate = ""
	c.Language = 0
	c.ImageQuality = "super_high"
	c.Split = "no"
}

func (c *Config) Validate() error {
	if _, ok := utils.FindInSliceFold([]string{"super_high", "high", "low"}, c.ImageQuality); !ok {
		return errors.New("InvalidMangaPlusImageQualityError: " + c.ImageQuality + " is not a valid image quality (super_high, high, low).")
	}
	if _, ok := utils.FindInSliceFold([]string{"no", "yes"}, c.Split); !ok {
		return errors.New("InvalidMangaPlusSplitError: " + c.Split + " is not a valid split value (no, yes).")
	}
	return nil
}
