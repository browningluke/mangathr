package register

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/database"
	"github.com/browningluke/mangathr/v2/internal/downloader"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/sources"
	"github.com/browningluke/mangathr/v2/internal/ui"
	"strings"
)

// Package-wide accessible driver
var driver *database.Driver

func closeDatabase() {
	err := driver.Close()
	if err != nil {
		logging.Errorln(err)
		ui.Error("Unable to close database.\nReason: %s\n", err)
	}
}

type options struct {
	title          string
	mapping        string
	includedGroups []string
	excludedGroups []string
	scraper        *sources.Scraper
}

func generateString(opts *options, prompt string) string {
	chapterTitles, err := (*opts.scraper).ChapterTitles()
	logging.ExitIfErrorWithFunc(err, closeDatabase)
	source := (*opts.scraper).ScraperName()

	latestChapter := "N/A"
	firstChapter := "N/A"

	if len(chapterTitles) > 0 {
		latestChapter = chapterTitles[0]
		firstChapter = chapterTitles[len(chapterTitles)-1]
	}

	return fmt.Sprintf(
		"\rTitle: %s"+
			"\nSource: %s"+
			"\n# of chapters: %d"+
			"\nLatest Chapter: %s"+
			"\nFirst  Chapter: %s"+
			"\nMapped to: ./%s"+
			"\nIncluded groups: [%s]"+
			"\nExcluded groups: [%s]"+
			"\n%s",
		opts.title, source, len(chapterTitles),
		latestChapter, firstChapter, opts.mapping,
		strings.Join(opts.includedGroups, ", "),
		strings.Join(opts.excludedGroups, ", "), prompt)
}

func findManga(args *registerOpts) (options, bool) {
	// Do manga scraping
	scraper := sources.NewScraper(args.Source)
	titles, err := scraper.Search(args.Query)
	logging.ExitIfErrorWithFunc(err, closeDatabase)

	// Check if scraper supports registering
	if !scraper.Registrable() {
		ui.PrintlnColor(ui.Yellow, "Selected scraper does not support registering to database. Exiting...")
		return options{}, true
	}

	// Exit if no manga are found
	if len(titles) == 0 {
		ui.PrintlnColor(ui.Yellow, "No manga found with specified query. Exiting...")
		return options{}, true
	}

	selection := titles[0]
	if len(titles) > 1 {
		var uierr error
		selection, uierr = ui.SingleCheckboxes("Select Manga:", titles)
		if uierr != nil {
			logging.ExitIfErrorWithFunc(&logging.ScraperError{
				Error: uierr, Message: "An error occurred while getting input", Code: 0,
			}, closeDatabase)
		}
	}

	err = scraper.SelectManga(selection)
	logging.ExitIfErrorWithFunc(err, closeDatabase)

	mangaTitle := scraper.MangaTitle()

	opts := options{
		title:          mangaTitle,
		mapping:        downloader.CleanPath(mangaTitle),
		scraper:        &scraper,
		includedGroups: []string{},
		excludedGroups: []string{},
	}

	if exists, _ := driver.CheckMangaExistence(scraper.MangaID()); exists {
		fmt.Printf("Title: %s\nSource: %s\n", mangaTitle, scraper.ScraperName())
		ui.PrintlnColor(ui.Yellow, "Manga is already registered. Exiting...")
		return options{}, true
	}

	// Manga might have some configured `filter` / `exclude` groups, so handle them here
	scraper.FilterGroups(opts.includedGroups, opts.excludedGroups)

	return opts, false
}

func handleMenu(args *registerOpts, driver *database.Driver) {
	opts, exists := findManga(args)
	if exists {
		return
	}

	// Define Customize panel
	customizePanel := ui.NewPanel().
		SetPrompt(func() string {
			return generateString(&opts, "Select an option")
		}).
		ErrorHandler(
			func(err error) {
				logging.ExitIfErrorWithFunc(&logging.ScraperError{
					Error: err, Message: "An error occurred while getting input", Code: 0,
				}, closeDatabase)
			},
		)

	customizePanel.
		AddOption("Change mapping").
		InputHandler("Map to:",
			func(i string) {
				// Handle input
				opts.mapping = downloader.CleanPath(i)
			},
			func(err error) {
				// Handle error
				logging.ExitIfErrorWithFunc(&logging.ScraperError{
					Error: err, Message: "An error occurred while getting input", Code: 0,
				}, closeDatabase)
			},
		)

	customizePanel.
		AddOption("Filter groups").
		CheckboxHandler("Select groups to filter: ",
			func() []string {
				// Generate options to display in checkboxes
				groups, err := (*opts.scraper).GroupNames()
				logging.ExitIfErrorWithFunc(err, closeDatabase)

				return groups
			},
			func(i []string) {
				// Handle selected options
				opts.includedGroups = i
				err := (*opts.scraper).FilterGroups(opts.includedGroups, opts.excludedGroups)
				logging.ExitIfErrorWithFunc(err, closeDatabase)
			},
			func(err error) {
				// Handle error
				logging.ExitIfErrorWithFunc(&logging.ScraperError{
					Error: err, Message: "An error occurred while getting input", Code: 0,
				}, closeDatabase)
			},
		)

	customizePanel.
		AddOption("Exclude groups").
		CheckboxHandler("Select groups to exclude: ",
			func() []string {
				// Generate options to display in checkboxes
				groups, err := (*opts.scraper).GroupNames()
				logging.ExitIfErrorWithFunc(err, closeDatabase)

				return groups
			},
			func(i []string) {
				// Handle selected options
				opts.excludedGroups = i
				err := (*opts.scraper).FilterGroups(opts.includedGroups, opts.excludedGroups)
				logging.ExitIfErrorWithFunc(err, closeDatabase)
			},
			func(err error) {
				// Handle error
				logging.ExitIfErrorWithFunc(&logging.ScraperError{
					Error: err, Message: "An error occurred while getting input", Code: 0,
				}, closeDatabase)
			},
		)

	// Define Main panel
	mainPanel := ui.NewPanel().
		SetPrompt(func() string {
			return generateString(&opts, "Select an option")
		}).
		ErrorHandler(
			func(err error) {
				logging.ExitIfErrorWithFunc(&logging.ScraperError{
					Error: err, Message: "An error occurred while getting input", Code: 0,
				}, closeDatabase)
			},
		)

	mainPanel.
		AddOption("Register").
		Terminator().
		ConfirmationHandler("Register this manga?",
			func() {
				// Handle yes
				// Start the registration process
				fmt.Println("confirmed")

				mangaID := (*opts.scraper).MangaID()
				source := (*opts.scraper).ScraperName()

				manga, err := driver.CreateManga(mangaID, opts.title, source, opts.mapping, opts.includedGroups, opts.excludedGroups)
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
			},
			func() {
				// Handle no
				fmt.Println("cancelled")
			},
			func(err error) {
				// Handle error
				logging.ExitIfErrorWithFunc(&logging.ScraperError{
					Error: err, Message: "An error occurred while getting input", Code: 0,
				}, closeDatabase)
			},
		)

	mainPanel.
		AddOption("Customize").
		PanelHandler(customizePanel)

	// This is a blocking call
	mainPanel.Start()
}
