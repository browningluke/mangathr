package mangaplus

var config Config

type Config struct {
	FilenameTemplate string `yaml:"filenameTemplate"`
	Language         int    `yaml:"language"`
}

func SetConfig(cfg Config) {
	config = cfg
}

func (c *Config) Default() {
	c.FilenameTemplate = ""
	c.Language = 0
}
