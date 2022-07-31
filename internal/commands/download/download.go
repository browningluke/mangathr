package download

import (
	"fmt"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/sources"
	"github.com/browningluke/mangathrV2/internal/ui"
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
	titles, err := scraper.Search(args.Query)
	logging.ExitIfError(err)

	selection := ui.SingleCheckboxes("Select Manga:", titles)
	err = scraper.SelectManga(selection)
	logging.ExitIfError(err)

	chapterTitles, err := scraper.ChapterTitles()
	logging.ExitIfError(err)

	//fmt.Println(chapters)
	chapterTitle := scraper.MangaTitle()
	sourceName := scraper.ScraperName()
	chapterSelections := SelectChapters(chapterTitles, chapterTitle, sourceName)
	//fmt.Println(chapterSelections)

	err = scraper.SelectChapters(chapterSelections)
	logging.ExitIfError(err)

	scraper.Download(downloader.NewDownloader(
		&config.Downloader, false,
		scraper.EnforceChapterDuration(),
	), "download")

}
