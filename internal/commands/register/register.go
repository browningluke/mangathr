package register

import (
	"fmt"
	"log"
	"mangathrV2/internal/config"
	"mangathrV2/internal/database"
	"mangathrV2/internal/sources"
	"mangathrV2/internal/utils/ui"
)

func SelectManga(titles []string) string {
	selection := ui.SingleCheckboxes(
		"Select Manga:",
		titles,
	)

	return selection
}

func ConfirmRegistration(titles []string, mangaTitle string, sourceName string) bool {
	return ui.ConfirmPrompt(
		fmt.Sprintf("\rTitle: %s\nSource: %s\n# of chapters: %d\nLatest Chapter: %s\nFirst  Chapter: %s\nRegister this manga?",
			mangaTitle, sourceName, len(titles), titles[0], titles[len(titles)-1]),
	)
}

func Run(args *Args, config *config.Config) {

	// Connect to database

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

	//_, err = driver.CreateManga("123414afsfaksjfgiaqa", "shikimori", "mangadex")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//_, err = driver.QueryManga("123414afsfaksjfgiaq")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	_, err = driver.QueryAllManga()
	if err != nil {
		log.Fatalln(err)
	}

	// Do manga scraping
	scraper := sources.NewScraper(args.Plugin, config)
	titles := scraper.Search(args.Query)
	selection := SelectManga(titles)
	scraper.SelectManga(selection)

	chapters := scraper.ListChapters()
	chapterTitle := scraper.GetMangaTitle()
	sourceName := scraper.GetScraperName()
	confirm := ConfirmRegistration(chapters, chapterTitle, sourceName)

	if confirm {
		fmt.Println("confirmed")
	} else {
		fmt.Println("cancelled")
	}

	// Prompt user to add to db (if not already in db)

	// Add to db

	// Close db

}
