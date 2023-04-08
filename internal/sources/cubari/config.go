package cubari

var config Config

type Config struct {
	FilenameTemplate string `yaml:"filenameTemplate"`
}

func SetConfig(cfg Config) {
	config = cfg
}

func (c *Config) Default() {
	c.FilenameTemplate = "" // No override of downloader.output.filenameTemplate
}
