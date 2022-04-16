package download

import (
	"fmt"
	"mangathrV2/internal/config"
	"mangathrV2/internal/downloader"
	"mangathrV2/internal/sources"
	"mangathrV2/internal/utils/ui"
)

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

	// Search and select manga
	titles := scraper.Search(args.Query)
	selection := ui.SingleCheckboxes("Select Manga:", titles)
	scraper.SelectManga(selection)

	chapterTitles := scraper.ChapterTitles()
	//fmt.Println(chapters)
	chapterTitle := scraper.MangaTitle()
	sourceName := scraper.ScraperName()
	chapterSelections := SelectChapters(chapterTitles, chapterTitle, sourceName)
	//fmt.Println(chapterSelections)
	scraper.SelectChapters(chapterSelections)

	scraper.Download(downloader.NewDownloader(&config.Downloader), "download")

}
