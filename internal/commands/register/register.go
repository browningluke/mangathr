package register

import (
	"fmt"
	"log"
	"mangathrV2/internal/config"
	"mangathrV2/internal/database"
	"mangathrV2/internal/sources"
	"mangathrV2/internal/utils/ui"
)

type options struct {
	title         string
	chapterTitles []string
	source        string
	mapping       string
	scraper       *sources.Scraper
}

func generateString(opts *options, prompt string) string {
	return fmt.Sprintf(
		"\rTitle: %s"+
			"\nSource: %s"+
			"\n# of chapters: %d"+
			"\nLatest Chapter: %s"+
			"\nFirst  Chapter: %s"+
			"\nMapped to: ./%s"+
			"\n%s",
		opts.title, opts.source, len(opts.chapterTitles), opts.chapterTitles[0],
		opts.chapterTitles[len(opts.chapterTitles)-1], opts.mapping, prompt)
}

func handleRegisterMenu(opts *options, driver *database.Driver) bool {
	confirm := ui.ConfirmPrompt("Register this manga?")
	if confirm {
		// Start the registration process
		fmt.Println("confirmed")

		mangaID := (*opts.scraper).MangaID()

		manga, err := driver.CreateManga(mangaID, opts.title, opts.source, opts.mapping)
		if err != nil {
			panic(err)
		}

		for _, c := range (*opts.scraper).Chapters() {
			err := driver.CreateChapter(c.ID, c.Num, manga)
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
	option := ui.SingleCheckboxes("Select an option", []string{"Change mapping", "Back"})

	switch option {
	case "Change mapping":
		res := ui.InputPrompt("Map to:")
		opts.mapping = res
		return true
	case "Back":
		return true
	default:
		panic("Option selected not in list")
	}

	return false
}

func promptMainMenu(args *Args, config *config.Config, driver *database.Driver) {
	// Do manga scraping
	scraper := sources.NewScraper(args.Plugin, config)
	titles := scraper.Search(args.Query)
	selection := ui.SingleCheckboxes("Select Manga:", titles)
	scraper.SelectManga(selection)

	chapterTitles := scraper.ChapterTitles()
	mangaTitle := scraper.MangaTitle()
	sourceName := scraper.ScraperName()
	opts := options{
		title:         mangaTitle,
		chapterTitles: chapterTitles,
		source:        sourceName,
		mapping:       mangaTitle,
		scraper:       &scraper,
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

func Run(args *Args, config *config.Config) {
	// Open database
	driver, err := database.GetDriver(database.SQLITE)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(driver *database.Driver) {
		err := driver.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(driver)

	promptMainMenu(args, config, driver)

	// Close database
}
