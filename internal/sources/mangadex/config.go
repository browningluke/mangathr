package mangadex

type Config struct {
	// Scraper
	FilenameTemplate string   `yaml:"filenameTemplate"`
	RatingFilter     []string `yaml:"ratingFilter"`
	LanguageFilter   []string `yaml:"languageFilter"`
	DataSaver        bool     `yaml:"dataSaver"`

	// Connection
	SyncDeletions bool `yaml:"syncDeletions"`
}
