package sources

import (
	"github.com/browningluke/mangathr/internal/downloader"
	"github.com/browningluke/mangathr/internal/logging"
	"github.com/browningluke/mangathr/internal/manga"
	"github.com/browningluke/mangathr/internal/sources/cubari"
	"github.com/browningluke/mangathr/internal/sources/mangadex"
	"github.com/browningluke/mangathr/internal/ui"
	"strings"
)

var scrapers = map[string]func() Scraper{
	// Mangadex
	strings.ToLower(mangadex.SCRAPERNAME): func() Scraper {
		return mangadex.NewScraper()
	},
	// Cubari
	strings.ToLower(cubari.SCRAPERNAME): func() Scraper {
		return cubari.NewScraper()
	},
}

var scraperTitles = map[string]string{
	// Mangadex
	strings.ToLower(mangadex.SCRAPERNAME): mangadex.SCRAPERNAME,
	// Cubari
	strings.ToLower(cubari.SCRAPERNAME): cubari.SCRAPERNAME,
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
	SelectNewChapters(chapterIDs []string) ([]manga.Chapter, *logging.ScraperError)
	SelectChapters(titles []string) *logging.ScraperError

	/*
		-- Chapter data --
	*/

	// Getters

	Chapters() ([]manga.Chapter, *logging.ScraperError)
	ChapterTitles() ([]string, *logging.ScraperError)
	GroupNames() ([]string, *logging.ScraperError)

	// Setters

	FilterGroups(groups []string) *logging.ScraperError

	// Downloading

	Download(downloader *downloader.Downloader, directoryMapping string) (succeeded []manga.Chapter)

	// Getters

	MangaTitle() string
	MangaID() string
	ScraperName() string
	EnforceChapterDuration() bool
	Registrable() bool
}

func MatchScraperTitle(query string) (string, bool) {
	matchedTitle, ok := scraperTitles[strings.ToLower(query)]
	return matchedTitle, ok
}

func NewScraper(name string) Scraper {
	getScraper, ok := scrapers[strings.ToLower(name)]
	if !ok {
		ui.Fatal("Scraper name could not be found.")
	}

	scraper := getScraper()
	logging.Infoln("Matched scraper: ", scraper.ScraperName())
	return scraper
}
