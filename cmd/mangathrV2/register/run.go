package register

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/database"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/sources"
	"github.com/browningluke/mangathrV2/internal/ui"
	"strings"
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

type options struct {
	title          string
	mapping        string
	filteredGroups []string
	scraper        *sources.Scraper
}

func generateString(opts *options, prompt string) string {
	chapterTitles, err := (*opts.scraper).ChapterTitles()
	logging.ExitIfErrorWithFunc(err, closeDatabase)
	source := (*opts.scraper).ScraperName()

	return fmt.Sprintf(
		"\rTitle: %s"+
			"\nSource: %s"+
			"\n# of chapters: %d"+
			"\nLatest Chapter: %s"+
			"\nFirst  Chapter: %s"+
			"\nMapped to: ./%s"+
			"\nFiltered groups: [%s]"+
			"\n%s",
		opts.title, source, len(chapterTitles), chapterTitles[0],
		chapterTitles[len(chapterTitles)-1], opts.mapping, strings.Join(opts.filteredGroups, ", "), prompt)
}

func handleRegisterMenu(opts *options, driver *database.Driver) bool {
	confirm := ui.ConfirmPrompt("Register this manga?")
	if confirm {
		// Start the registration process
		fmt.Println("confirmed")

		mangaID := (*opts.scraper).MangaID()
		source := (*opts.scraper).ScraperName()

		manga, err := driver.CreateManga(mangaID, opts.title, source, opts.mapping, opts.filteredGroups)
		if err != nil {
			logging.ExitIfErrorWithFunc(&logging.ScraperError{
				Error: err, Message: "An error occurred when adding Manga to database", Code: 0,
			}, closeDatabase)
		}

		chapters, scraperErr := (*opts.scraper).Chapters()
		logging.ExitIfErrorWithFunc(scraperErr, closeDatabase)

		for _, c := range chapters {
			err := driver.CreateChapter(c.ID, c.Metadata.Num, c.Metadata.Title, manga)
			if err != nil {
				panic(err)
			}
		}

	} else {
		// return
		fmt.Println("cancelled")
	}

	return false
}

func handleCustomizeMenu(opts *options) bool {
	option := ui.SingleCheckboxes("Select an option",
		[]string{"Change mapping", "Filter groups", "Back"})

	switch option {
	case "Change mapping":
		res := ui.InputPrompt("Map to:")
		opts.mapping = downloader.CleanPath(res)
		return true
	case "Filter groups":
		groups, err := (*opts.scraper).GroupNames()
		logging.ExitIfErrorWithFunc(err, closeDatabase)

		selectedGroups := ui.Checkboxes("Select groups to filter: ", groups)
		opts.filteredGroups = selectedGroups

		err = (*opts.scraper).FilterGroups(selectedGroups)
		logging.ExitIfErrorWithFunc(err, closeDatabase)
		return true
	case "Back":
		return true
	default:
		panic("Option selected not in list")
	}

	return false
}

func promptMainMenu(args *registerOpts, config *config.Config, driver *database.Driver) {
	// Do manga scraping
	scraper := sources.NewScraper(args.Source, config)
	titles, err := scraper.Search(args.Query)
	logging.ExitIfErrorWithFunc(err, closeDatabase)

	selection := titles[0]
	if len(titles) > 1 {
		selection = ui.SingleCheckboxes("Select Manga:", titles)
	}

	err = scraper.SelectManga(selection)
	logging.ExitIfErrorWithFunc(err, closeDatabase)

	mangaTitle := scraper.MangaTitle()

	opts := options{
		title:          mangaTitle,
		mapping:        downloader.CleanPath(mangaTitle),
		scraper:        &scraper,
		filteredGroups: []string{},
	}

	if exists, _ := driver.CheckMangaExistence(scraper.MangaID()); exists {
		fmt.Printf("Title: %s\nSource: %s\n", mangaTitle, scraper.ScraperName())
		ui.PrintlnColor(ui.Yellow, fmt.Sprint("Manga is already registered. Exiting..."))
		return
	}

	for true {
		option := ui.SingleCheckboxes(generateString(&opts, "Select an option"),
			[]string{"Register", "Customize"})

		if option == "Register" {
			if loop := handleRegisterMenu(&opts, driver); !loop {
				break
			}
		} else if option == "Customize" {
			if loop := handleCustomizeMenu(&opts); !loop {
				break
			}
		}
	}
}
