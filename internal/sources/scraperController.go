package sources

import (
	"mangathrV2/internal/config"
	"mangathrV2/internal/downloader"
	"mangathrV2/internal/sources/mangadex"
)

type Scraper interface {
	Search(query string) []string
	SearchByID(id string) interface{}

	SelectManga(name string)
	ListChapters() []string

	//SelectNewChapters() interface{}
	SelectChapters(titles []string)

	Download(downloader *downloader.Downloader, downloadType string)

	GetMangaTitle() string
	GetScraperName() string
}

func NewScraper(name string, config *config.Config) Scraper {
	m := map[string]func() Scraper{
		"mangadex": func() Scraper { return mangadex.NewScraper(&config.Sources.Mangadex) },
		//"cubari":   func() Scraper { return cubari.NewScraper() },
	}

	scraper, ok := m[name]
	if !ok {
		panic("Passed scraper name not in map")
	}
	return scraper()
}
