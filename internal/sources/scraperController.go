package sources

import (
	"mangathrV2/internal/config"
	"mangathrV2/internal/downloader"
	"mangathrV2/internal/sources/mangadex"
	"mangathrV2/internal/sources/structs"
)

type Scraper interface {
	// Searching

	Search(query string) []string
	SearchByID(id string) interface{}

	// Selection

	SelectManga(name string)
	//SelectNewChapters() interface{}
	SelectChapters(titles []string)

	// Downloading

	Download(downloader *downloader.Downloader, downloadType string)

	// Getters

	MangaID() string
	MangaTitle() string

	Chapters() []structs.Chapter
	ChapterTitles() []string

	ScraperName() string

	EnforceChapterLength() bool
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
