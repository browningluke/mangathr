package downloader

import (
	"log"
	"os"
	"path/filepath"
)

type Downloader struct {
	config *Config
}

func NewDownloader(config *Config) *Downloader {
	return &Downloader{config: config}
}

func (d *Downloader) CreateDirectory(title string, downloadType string) {
	var dirname string

	if downloadType == "download" {
		if d.config.Output.Path == "" {
			d, err := os.UserHomeDir()
			if err != nil {
				log.Fatalln(err)
			}
			dirname = filepath.Join(d, "mangathrV2")
		} else {
			dirname = d.config.Output.Path
		}
	} else {
		if d.config.Output.UpdatePath == "" {
			d, err := os.UserHomeDir()
			if err != nil {
				log.Fatalln(err)
			}
			dirname = filepath.Join(d, "mangathrV2")
		} else {
			dirname = d.config.Output.UpdatePath
		}
	}

	newPath := filepath.Join(dirname, title)
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
}
