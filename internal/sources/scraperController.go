package sources

import (
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/sources/mangadex"
	"github.com/browningluke/mangathrV2/internal/sources/structs"
	"strings"
)

var SCRAPERS = map[string]func(c *config.Config) Scraper{
	// Mangadex
	strings.ToLower(mangadex.SCRAPERNAME): func(c *config.Config) Scraper {
		return mangadex.NewScraper(&c.Sources.Mangadex)
	},
	// Cubari
	//strings.ToLower(cubari.SCRAPERNAME): func(c *config.Config) Scraper {
	//	return cubari.NewScraper(&c.Sources.Cubari)
	//},
}

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
	getScraper, ok := SCRAPERS[strings.ToLower(name)]
	if !ok {
		panic("Passed scraper name not in map")
	}
	return getScraper(config)
}
