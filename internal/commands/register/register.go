package register

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/database"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/sources"
	"github.com/browningluke/mangathrV2/internal/utils/ui"
	"log"
	"strings"
)

type options struct {
	title          string
	mapping        string
	filteredGroups []string
	scraper        *sources.Scraper
}

func generateString(opts *options, prompt string) string {
	chapterTitles := (*opts.scraper).ChapterTitles()
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
			panic(err)
		}

		for _, c := range (*opts.scraper).Chapters() {
			err := driver.CreateChapter(c.ID, c.Num, c.Title, manga)
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
		groups := (*opts.scraper).GroupNames()
		selectedGroups := ui.Checkboxes("Select groups to filter: ", groups)
		opts.filteredGroups = selectedGroups
		(*opts.scraper).FilterGroups(selectedGroups)
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

	selection := titles[0]
	if len(titles) > 1 {
		selection = ui.SingleCheckboxes("Select Manga:", titles)
	}
	scraper.SelectManga(selection)
	mangaTitle := scraper.MangaTitle()

	opts := options{
		title:          mangaTitle,
		mapping:        downloader.CleanPath(mangaTitle),
		scraper:        &scraper,
		filteredGroups: []string{},
	}

	if exists, _ := driver.CheckMangaExistence(scraper.MangaID()); exists {
		fmt.Printf("Title: %s\nSource: %s\n", mangaTitle, scraper.ScraperName())
		ui.PrintlnColor(fmt.Sprint("Manga is already registered. Exiting..."), ui.Red)
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

func Run(args *Args, config *config.Config) {
	// Open database
	driver, err := database.GetDriver(database.SQLITE, config.Database.Uri)
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
}
