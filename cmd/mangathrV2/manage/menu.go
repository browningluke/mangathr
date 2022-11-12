package manage

import (
	"fmt"
	"github.com/browningluke/mangathrV2/ent"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/database"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/ui"
	"github.com/browningluke/mangathrV2/internal/utils"
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

func deleteFromDatabase(filter func(manga *ent.Manga) bool) {
	allManga, err := driver.QueryAllManga()
	if err != nil {
		logging.ExitIfErrorWithFunc(&logging.ScraperError{
			Error: err, Message: "An error occurred while getting manga from database", Code: 0,
		}, closeDatabase)
	}

	for _, manga := range allManga {
		if filter(manga) {
			err := driver.DeleteManga(manga)
			if err != nil {
				panic(err)
			}
		}
	}
}

func handleMenu(args *manageOpts, config *config.Config, driver *database.Driver) {
	// Define Main panel
	mainPanel := ui.NewPanel().
		SetPrompt(func() string {
			return "Select an option"
		}).
		ErrorHandler(
			func(err error) {
				logging.ExitIfErrorWithFunc(&logging.ScraperError{
					Error: err, Message: "An error occurred while getting input", Code: 0,
				}, closeDatabase)
			},
		)

	mainPanel.
		AddOption("Delete").
		Terminator().
		CheckboxHandler("Select manga to delete",
			func() []string {
				allManga, err := driver.QueryAllManga()
				if err != nil {
					logging.ExitIfErrorWithFunc(&logging.ScraperError{
						Error: err, Message: "An error occurred while getting manga from database", Code: 0,
					}, closeDatabase)
				}

				var mangaSelections []string
				for _, manga := range allManga {
					mangaSelections = append(mangaSelections, fmt.Sprintf("[%s] %s", manga.Source, manga.Title))
				}

				return mangaSelections
			},
			func(strings []string) {
				// Do nothing if selection list is empty
				if len(strings) == 0 {
					return
				}

				deleteFromDatabase(func(manga *ent.Manga) bool {
					_, exists := utils.FindInSlice(strings,
						fmt.Sprintf("[%s] %s", manga.Source, manga.Title))
					return exists
				})

				fmt.Println("Successfully deleted selected manga")
			},
			func(err error) {
				if err.Error() == "please provide options to select from" {
					fmt.Println("No manga are registered")
					return
				}

				logging.ExitIfErrorWithFunc(&logging.ScraperError{
					Error: err, Message: "An error occurred while getting input", Code: 0,
				}, closeDatabase)
			},
		)

	mainPanel.Start()

}
