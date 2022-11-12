package manage

import (
	"fmt"
	"github.com/browningluke/mangathrV2/ent"
	"github.com/browningluke/mangathrV2/ent/manga"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/database"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/sources"
	"github.com/browningluke/mangathrV2/internal/ui"
	"strings"
)

func deleteFromDatabase(filter func(manga *ent.Manga) bool) {
	// todo: use ent search, rather than querying all
	allManga, err := driver.QueryAllManga()
	if err != nil {
		logging.ExitIfErrorWithFunc(&logging.ScraperError{
			Error: err, Message: "An error occurred while getting manga from database", Code: 0,
		}, closeDatabase)
	}

	for _, m := range allManga {
		if filter(m) {
			err := driver.DeleteManga(m)
			if err != nil {
				panic(err)
			}
		}
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
		deleteFromDatabase(func(manga *ent.Manga) bool {
			return strings.ToLower(manga.Title) == strings.ToLower(seriesTitle) &&
				manga.Source == sourceTitle
		})

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
