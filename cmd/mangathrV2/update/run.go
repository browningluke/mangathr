package update

import (
	"fmt"
	"github.com/browningluke/mangathrV2/ent"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/database"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/sources"
	"github.com/browningluke/mangathrV2/internal/ui"
)

// Package-wide accessible driver
var driver *database.Driver

func closeDatabase() {
	logging.Warningln("Closing database because of error")
	err := driver.Close()
	if err != nil {
		logging.Errorln(err)
		ui.Error("Unable to close database.")
	}
}

func downloadNewChapters(config *config.Config, manga *ent.Manga, scraper sources.Scraper, numChapters int) {
	fmt.Printf("\033[2K") // Clear line
	fmt.Printf(fmt.Sprintf("\rTitle: %s\nSource: %s\n# of new chapters: %d\n",
		scraper.MangaTitle(), scraper.ScraperName(), numChapters))

	succeeded := scraper.Download(downloader.NewDownloader(
		&config.Downloader, true,
		scraper.EnforceChapterDuration()), "update")

	if !config.Downloader.DryRun {
		// If it's not a dry run, add new chapters to db
		logging.Debugln("Saving chapters to db")

		// Loop through successfully downloaded chapters, and add them to the db
		// (will retry failed chapters on next run)
		for _, chapter := range succeeded {
			err := driver.CreateChapter(chapter.ID, chapter.Metadata.Num, chapter.Metadata.Title, manga)
			if err != nil {
				ui.Error("Failed to save chapter to db: ",
					chapter.Metadata.Title, " (", chapter.ID, ")")
			}
		}
	}
}

func checkMangaForNewChapters(config *config.Config, manga *ent.Manga) {
	logging.Debugln("Requesting source...", manga.Source)
	scraper := sources.NewScraper(manga.Source, config)

	// Directly search for chapter by ID
	if err := scraper.SearchByID(manga.MangaID, manga.Title); err != nil {
		// Log error, abandon search, and continue (no need to exit program)
		logging.Errorln(err)
		ui.Error("An error occurred while search for ", manga.Title)
	}

	fmt.Printf("Checking  %s", manga.Title)

	// Convert ent chapters to chapterID array
	var chapterIDs []string
	for _, chapter := range manga.Edges.Chapters {
		chapterIDs = append(chapterIDs, chapter.ChapterID)
	}

	// Select new chapters in scraper, get array of them; and download if > 0
	newChapters, err := scraper.SelectNewChapters(chapterIDs)
	if err != nil {
		// Log error, abandon search, and continue (no need to exit program)
		logging.Errorln(err)
		ui.Error("An error occurred while search for ", manga.Title)
	}

	if numChapters := len(newChapters); numChapters > 0 {
		downloadNewChapters(config, manga, scraper, numChapters)
	} else {
		fmt.Printf("\rNone for  %s\n", manga.Title)
	}
}