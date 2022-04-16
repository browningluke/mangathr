package register

import (
	"fmt"
	"mangathrV2/internal/config"
	"mangathrV2/internal/sources"
	"mangathrV2/internal/utils/ui"
)

type options struct {
	title    string
	chapters []string
	source   string
	mapping  string
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
		opts.title, opts.source, len(opts.chapters), opts.chapters[0],
		opts.chapters[len(opts.chapters)-1], opts.mapping, prompt)
}

func handleRegisterMenu(opts *options) bool {
	confirm := ui.ConfirmPrompt("Register this manga?")
	if confirm {
		// Start the registration process
		fmt.Println("confirmed")
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

func promptMainMenu(args *Args, config *config.Config) {
	// Do manga scraping
	scraper := sources.NewScraper(args.Plugin, config)
	titles := scraper.Search(args.Query)
	selection := ui.SingleCheckboxes("Select Manga:", titles)
	scraper.SelectManga(selection)

	chapters := scraper.ListChapters()
	mangaTitle := scraper.GetMangaTitle()
	sourceName := scraper.GetScraperName()
	opts := options{title: mangaTitle, chapters: chapters, source: sourceName, mapping: mangaTitle}

	for true {
		option := ui.SingleCheckboxes(generateString(&opts, "Select an option"),
			[]string{"Register", "Customize"})

		if option == "Register" {
			if loop := handleRegisterMenu(&opts); !loop {
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

	promptMainMenu(args, config)

	// Close database
}
