package update

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/ent"
	"github.com/browningluke/mangathr/v2/internal/database"
	"github.com/browningluke/mangathr/v2/internal/downloader"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/sources"
	"github.com/browningluke/mangathr/v2/internal/ui"
)

// Package-wide accessible driver
var driver *database.Driver

func closeDatabase() {
	err := driver.Close()
	if err != nil {
		logging.Errorln(err)
		ui.Errorf("Unable to close database.\nReason: %s\n", err)
	}
}

func downloadNewChapters(manga *ent.Manga,
	scraper sources.Scraper, numChapters int) (downloaded, errors int) {

	fmt.Printf("\033[2K") // Clear line
	fmt.Printf(
		"\r\u001B[1mTitle:  \u001B[0m%s\n"+
			"\u001B[1mSource: \u001B[0m%s\n"+
			"\u001B[1mFound:  \u001B[0m%d chapter(s)\n",
		scraper.MangaTitle(), scraper.ScraperName(), numChapters)

	succeeded := scraper.Download(
		downloader.NewDownloader(downloader.UPDATE, scraper.EnforceChapterDuration()), manga.Mapping,
	)

	if !downloader.DryRun() {
		// If it's not a dry run, add new chapters to db
		logging.Debugln("Saving chapters to db")

		// Loop through successfully downloaded chapters, and add them to the db
		// (will retry failed chapters on next run)
		for _, chapter := range succeeded {
			err := driver.CreateChapter(chapter.ID, chapter.Metadata.Num, chapter.Metadata.Title, manga)
			if err != nil {
				logging.Errorln(err)
				ui.Error("Failed to save chapter to db: ",
					chapter.Metadata.Title, " (", chapter.ID, ")")
			}
		}
	}

	return len(succeeded), numChapters - len(succeeded)
}

func checkMangaForNewChapters(manga *ent.Manga) (seriesStats, *logging.ScraperError) {
	stats := seriesStats{}

	logging.Debugln("Requesting source...", manga.Source)
	scraper := sources.NewScraper(manga.Source)

	// Directly search for chapter by ID
	if err := scraper.SearchByID(manga.MangaID, manga.Title); err != nil {
		// Abandon search
		return seriesStats{}, err
	}

	fmt.Printf("Checking  %s", manga.Title)

	// Convert ent chapters to chapterID array
	var chapterIDs []string
	for _, chapter := range manga.Edges.Chapters {
		chapterIDs = append(chapterIDs, chapter.ChapterID)
	}

	// Filter groups
	scraper.FilterGroups(manga.FilteredGroups, manga.ExcludedGroups)

	// Select new chapters in scraper, get array of them; and download if > 0
	newChapters, err := scraper.SelectNewChapters(chapterIDs)
	if err != nil {
		// Log error, abandon search, and continue (no need to exit program)
		logging.Errorln(err)
		ui.Error("An error occurred while search for ", manga.Title)
	}

	stats.found = len(newChapters)

	if numChapters := len(newChapters); numChapters > 0 {
		stats.downloaded, stats.errors = downloadNewChapters(manga, scraper, numChapters)
	} else {
		fmt.Printf("\rNone for  %s\n", manga.Title)
	}

	return stats, nil
}
