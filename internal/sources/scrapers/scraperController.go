package scrapers

import (
	"mangathrV2/internal/downloader"
	"mangathrV2/internal/sources/scrapers/mangadex"
)

type Scraper interface {
	Search(query string) []string
	SearchByID(id string) interface{}

	SelectManga(name string)
	ListChapters() []string

	//SelectNewChapters() interface{}
	SelectChapters(titles []string)

	Download(downloader *downloader.Downloader)

	GetChapterTitle() string
	GetScraperName() string
}

func NewScraper(name string) Scraper {
	m := map[string]func() Scraper{
		"mangadex": func() Scraper { return mangadex.NewScraper() },
		//"cubari":   func() Scraper { return cubari.NewScraper() },
	}

	scraper, ok := m[name]
	if !ok {
		panic("Passed scraper name not in map")
	}
	return scraper()
}
