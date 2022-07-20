package sources

import (
	"mangathrV2/internal/config"
	"mangathrV2/internal/downloader"
	"mangathrV2/internal/sources/mangadex"
	"mangathrV2/internal/sources/structs"
	"strings"
)

type Scraper interface {
	/*
		-- Searching --
	*/

	Search(query string) []string
	SearchByID(id, title string) error

	/*
		-- Chapter scraping --
	*/

	SelectManga(name string)
	SelectNewChapters(chapters []structs.Chapter) []structs.Chapter
	SelectChapters(titles []string)

	/*
		-- Chapter data --
	*/

	// Getters

	Chapters() []structs.Chapter
	ChapterTitles() []string
	GroupNames() []string

	// Setters

	FilterGroups(groups []string)

	// Downloading

	Download(downloader *downloader.Downloader, downloadType string)

	// Getters

	MangaTitle() string
	MangaID() string
	ScraperName() string
	EnforceChapterDuration() bool
}

func NewScraper(name string, config *config.Config) Scraper {
	name = strings.ToLower(name)

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
