package download

import (
	"fmt"
	"mangathrV2/internal/config"
	"mangathrV2/internal/downloader"
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

func SelectChapters(titles []string, mangaTitle string, sourceName string) []string {

	selections := ui.Checkboxes(
		fmt.Sprintf("\rTitle: %s\nSource: %s\n# of chapters: %d\nSelect chapters",
			mangaTitle, sourceName, len(titles)),
		titles,
	)

	return selections
}

func Run(args *Args, config *config.Config) {
	scraper := sources.NewScraper(args.Plugin, config)

	titles := scraper.Search(args.Query)
	//fmt.Println(titles)

	selection := SelectManga(titles)
	scraper.SelectManga(selection)

	chapters := scraper.ListChapters()
	//fmt.Println(chapters)
	chapterTitle := scraper.GetMangaTitle()
	sourceName := scraper.GetScraperName()
	chapterSelections := SelectChapters(chapters, chapterTitle, sourceName)
	//fmt.Println(chapterSelections)
	scraper.SelectChapters(chapterSelections)

	scraper.Download(downloader.NewDownloader(&config.Downloader), "download")

}
