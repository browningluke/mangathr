package manage

import (
	"fmt"
	"github.com/browningluke/mangathr/ent"
	"github.com/browningluke/mangathr/ent/manga"
	"github.com/browningluke/mangathr/internal/config"
	"github.com/browningluke/mangathr/internal/database"
	"github.com/browningluke/mangathr/internal/logging"
	"github.com/browningluke/mangathr/internal/sources"
	"github.com/browningluke/mangathr/internal/ui"
)

// deleteFromDatabase removes manga by source + query
func deleteFromDatabase(source, query string, isTitle bool) {
	var queriedManga *ent.Manga
	var err error

	if isTitle {
		queriedManga, err = driver.QueryMangaByTitle(query, source)
	} else {
		queriedManga, err = driver.QueryMangaByID(query, source)
	}

	if err != nil {
		logging.ExitIfErrorWithFunc(&logging.ScraperError{
			Error: err, Message: "An error occurred while getting manga from database", Code: 0,
		}, closeDatabase)
	}

	err = driver.DeleteManga(queriedManga)
	if err != nil {
		logging.ExitIfErrorWithFunc(&logging.ScraperError{
			Error: err, Message: "An error occurred while deleting manga from database", Code: 0,
		}, closeDatabase)
	}
}

func promptDeletion(sourceTitle, seriesTitle string) {
	fmt.Printf("Deleting [%s] %s\n", sourceTitle, seriesTitle)
	confirm, err := ui.ConfirmPrompt("Delete this manga?")
	if err != nil {
		logging.ExitIfErrorWithFunc(&logging.ScraperError{
			Error: err, Message: "An error occurred while getting input", Code: 0,
		}, closeDatabase)
	}

	if confirm {
		fmt.Println("Deleting series")
		deleteFromDatabase(sourceTitle, seriesTitle, true)
	} else {
		fmt.Println("Skipping deletion...")
	}
}

func handleDelete(opts *manageOpts, c *config.Config, d *database.Driver) {

	for _, title := range opts.Delete.Query {
		// Ignoring ok value, since check was performed before this function was called
		sourceTitle, _ := sources.MatchScraperTitle(opts.Delete.Source)

		exists, _ := d.CheckMangaExistenceByPredicate(
			manga.TitleEqualFold(title),
			manga.Source(sourceTitle),
		)

		if !exists {
			ui.PrintfColor(ui.Red, "No series with title \"%s\" found\n", title)
			continue
		}
		promptDeletion(sourceTitle, title)
	}

}
