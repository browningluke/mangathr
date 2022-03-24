package scrapers

import (
	"mangathrV2/internal/downloader"
	"mangathrV2/internal/sources/scrapers/cubari"
	"mangathrV2/internal/sources/scrapers/mangadex"
)

type Scraper interface {
	Search(query string) interface{}
	SearchByID(id string) interface{}
	ListChapters() interface{}

	//SelectNewChapters() interface{}
	SelectChapters() interface{}

	Download(downloader *downloader.Downloader)
}

func NewScraper(name string) Scraper {
	m := map[string]func() Scraper{
		"mangadex": func() Scraper { return mangadex.NewScraper() },
		"cubari":   func() Scraper { return cubari.NewScraper() },
	}

	scraper, ok := m[name]
	if !ok {
		panic("Passed scraper name not in map")
	}
	return scraper()
}
