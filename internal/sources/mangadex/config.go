package mangadex

var config Config

type Config struct {
	// Scraper
	FilenameTemplate string   `yaml:"filenameTemplate"`
	RatingFilter     []string `yaml:"ratingFilter"`
	LanguageFilter   []string `yaml:"languageFilter"`
	DataSaver        bool     `yaml:"dataSaver"`

	// Groups
	Groups struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"groups"`

	// Connection
	SyncDeletions bool `yaml:"syncDeletions"`
}

func SetConfig(cfg Config) {
	config = cfg
}

func (c *Config) Default() {
	c.FilenameTemplate = "" // No override of downloader.output.filenameTemplate
	c.RatingFilter = []string{"safe", "suggestive", "erotica"}
	c.LanguageFilter = []string{"en"}
	c.DataSaver = false

	c.Groups.Include = []string{}
	c.Groups.Exclude = []string{}

	c.SyncDeletions = false
}
