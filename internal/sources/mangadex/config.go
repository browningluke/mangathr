package mangadex

type Config struct {
	// Scraper
	DataSaver bool `yaml:"dataSaver"`

	// Connection
	SyncDeletions bool `yaml:"syncDeletions"`
}
