package downloader

type Config struct {
	SimultaneousPages int `yaml:"simultaneousPages"`
	PageRetries       int `yaml:"pageRetries"`
	Delay             struct {
		Page          int
		Chapter       int
		UpdateChapter int `yaml:"updateChapter"`
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
