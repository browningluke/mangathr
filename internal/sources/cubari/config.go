package cubari

var config Config

type Config struct {
	// Scraper
	FilenameTemplate string `yaml:"filenameTemplate"`

	// Groups
	Groups struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"groups"`
}

func SetConfig(cfg Config) {
	config = cfg
}

func (c *Config) Default() {
	c.FilenameTemplate = "" // No override of downloader.output.filenameTemplate
	c.Groups.Include = []string{}
	c.Groups.Exclude = []string{}
}
