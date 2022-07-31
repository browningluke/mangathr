package sources

import (
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/sources/mangadex"
	"github.com/browningluke/mangathrV2/internal/sources/structs"
	"github.com/browningluke/mangathrV2/internal/ui"
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

	Search(query string) ([]string, *logging.ScraperError)
	SearchByID(id, title string) *logging.ScraperError

	/*
		-- Chapter scraping --
	*/

	SelectManga(name string) *logging.ScraperError
	SelectNewChapters(chapters []structs.Chapter) ([]structs.Chapter, *logging.ScraperError)
	SelectChapters(titles []string) *logging.ScraperError

	/*
		-- Chapter data --
	*/

	// Getters

	Chapters() ([]structs.Chapter, *logging.ScraperError)
	ChapterTitles() ([]string, *logging.ScraperError)
	GroupNames() ([]string, *logging.ScraperError)

	// Setters

	FilterGroups(groups []string) *logging.ScraperError

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
		ui.Fatal("Scraper name could not be found.")
	}

	scraper := getScraper(config)
	logging.Infoln("Matched scraper: ", scraper.ScraperName())
	return scraper
}
