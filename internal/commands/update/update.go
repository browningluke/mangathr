package update

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/database"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/sources"
	"github.com/browningluke/mangathrV2/internal/sources/structs"
	"github.com/browningluke/mangathrV2/internal/ui"
	"time"
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

func Run(config *config.Config) {
	// Open database
	var err error
	driver, err = database.GetDriver(database.SQLITE, config.Database.Uri)
	if err != nil {
		logging.Errorln(err)
		ui.Fatal("Unable to open database.")
	}
	defer closeDatabase()

	allManga, err := driver.QueryAllManga()
	if err != nil {
		logging.ExitIfErrorWithFunc(&logging.ScraperError{
			Error: err, Message: "An error occurred while getting Manga from database.", Code: 0,
		}, closeDatabase)
	}

	for _, manga := range allManga {
		logging.Debugln("Requesting source...", manga.Source)
		scraper := sources.NewScraper(manga.Source, config)

		if err := scraper.SearchByID(manga.MangaID, manga.Title); err != nil {
			// Log error, abandon search, and continue (no need to exit program)
			logging.Errorln(err)
			ui.Error("An error occurred while search for ", manga.Title)
		}

		fmt.Printf("Checking  %s", manga.Title)

		// Convert ent chapters to chapter struct
		var chapters []structs.Chapter
		for _, chapter := range manga.Edges.Chapters {
			chapters = append(chapters, structs.Chapter{
				ID:    chapter.ChapterID,
				Num:   chapter.Num,
				Title: chapter.Title,
				Metadata: structs.Metadata{
					Date:   "",  // These can be empty since
					Link:   "",  // these will not be downloaded
					Groups: nil, // and metadata is not needed.
				},
			})
		}

		newChapters, err := scraper.SelectNewChapters(chapters)
		if err != nil {
			// Log error, abandon search, and continue (no need to exit program)
			logging.Errorln(err)
			ui.Error("An error occurred while search for ", manga.Title)
		}

		if len(newChapters) > 0 {
			fmt.Printf("\033[2K") // Clear line
			fmt.Printf(fmt.Sprintf("\rTitle: %s\nSource: %s\n# of new chapters: %d\n",
				scraper.MangaTitle(), scraper.ScraperName(), len(newChapters)))

			scraper.Download(downloader.NewDownloader(
				&config.Downloader, true,
				scraper.EnforceChapterDuration()), "update")

			if !config.Downloader.DryRun {
				// update in db
				logging.Debugln("Saving chapters to db")

				// todo. here we assume all downloads succeeded.
				// todo. figure out how to determine which downloads failed
				for _, chapter := range newChapters {
					err := driver.CreateChapter(chapter.ID, chapter.Num, chapter.Title, manga)
					if err != nil {
						ui.Error("Failed to save chapter to db: ",
							chapter.Title, " (", chapter.ID, ")")
					}
				}
			}
		} else {
			fmt.Printf("\rNone for  %s\n", manga.Title)
		}

		// Sleep between checks
		dur, durErr := time.ParseDuration(config.Downloader.Delay.UpdateChapter)
		if durErr != nil {
			logging.Errorln(durErr)
			ui.Error("Failed to parse time duration.")
		}
		time.Sleep(dur)

	}
}
