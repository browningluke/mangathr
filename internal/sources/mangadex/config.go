package mangadex

var config Config

type Config struct {
	// Scraper
	FilenameTemplate string   `yaml:"filenameTemplate"`
	RatingFilter     []string `yaml:"ratingFilter"`
	LanguageFilter   []string `yaml:"languageFilter"`
	DataSaver        bool     `yaml:"dataSaver"`

	// Connection
	SyncDeletions bool `yaml:"syncDeletions"`
}

func SetConfig(cfg Config) {
	config = cfg
}

func (c *Config) Default() {
	c.FilenameTemplate = "" // No override of downloader.output.filenameTemplate
	c.RatingFilter = []string{"safe", "suggestive"}
	c.LanguageFilter = []string{"en"}
	c.DataSaver = false

	c.SyncDeletions = false
}
