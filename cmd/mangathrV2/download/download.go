package download

import (
	"errors"
	"fmt"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/downloader"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/sources"
	"github.com/browningluke/mangathrV2/internal/ui"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *config.Config) *cobra.Command {
	o := &downloadOpts{}

	cmd := &cobra.Command{
		Use:     "download [OPTIONS] -s SOURCE QUERY",
		Short:   "Download chapters from source",
		Aliases: []string{"d"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a query to search")
			}
			if len(args) > 1 {
				return errors.New("can only search 1 query at a time")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Query = args[0]
			cfg.Downloader.DryRun = cfg.Downloader.DryRun || o.DryRun
			o.run(cfg)
		},
		DisableFlagsInUseLine: true,
	}

	cmd.Flags().StringVarP(&o.Source, "source", "s",
		"", "source to search")
	err := cmd.MarkFlagRequired("source")
	cobra.CheckErr(err)

	cmd.Flags().BoolVarP(&o.DryRun, "dry-run", "",
		false, "do not download files or update database")

	return cmd
}

func (o *downloadOpts) run(cfg *config.Config) {
	scraper := sources.NewScraper(o.Source, cfg)

	// Search and select manga
	titles, err := scraper.Search(o.Query)
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
		&cfg.Downloader, false,
		scraper.EnforceChapterDuration(),
	), "download")
}

func SelectChapters(titles []string, mangaTitle string, sourceName string) []string {

	selections := ui.Checkboxes(
		fmt.Sprintf("\rTitle: %s\nSource: %s\n# of chapters: %d\nSelect chapters",
			mangaTitle, sourceName, len(titles)),
		titles,
	)

	return selections
}
