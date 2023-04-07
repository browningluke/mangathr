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
		"", "Source to search query on")
	err := cmd.MarkFlagRequired("source")
	cobra.CheckErr(err)

	cmd.Flags().BoolVarP(&o.DryRun, "dry-run", "",
		false, "Disables downloads & writes to the database")

	return cmd
}

func (o *downloadOpts) run(cfg *config.Config) {
	scraper := sources.NewScraper(o.Source, cfg)
	cfg.Propagate()

	// Search and select manga
	titles, err := scraper.Search(o.Query)
	logging.ExitIfError(err)

	selection, uierr := ui.SingleCheckboxes("Select Manga:", titles)
	if uierr != nil {
		logging.ExitIfError(&logging.ScraperError{
			Error: uierr, Message: "An error occurred while getting input", Code: 0,
		})
	}

	err = scraper.SelectManga(selection)
	logging.ExitIfError(err)

	chapterTitles, err := scraper.ChapterTitles()
	logging.ExitIfError(err)

	//fmt.Println(chapters)
	chapterTitle := scraper.MangaTitle()
	sourceName := scraper.ScraperName()
	chapterSelections, uierr := SelectChapters(chapterTitles, chapterTitle, sourceName)
	if uierr != nil {
		logging.ExitIfError(&logging.ScraperError{
			Error: uierr, Message: "An error occurred while getting input", Code: 0,
		})
	}

	err = scraper.SelectChapters(chapterSelections)
	logging.ExitIfError(err)

	scraper.Download(
		downloader.NewDownloader(false, scraper.EnforceChapterDuration()),
		"", "download")
}

func SelectChapters(titles []string, mangaTitle string, sourceName string) ([]string, error) {

	selections, err := ui.Checkboxes(
		fmt.Sprintf("\rTitle: %s\nSource: %s\n# of chapters: %d\nSelect chapters",
			mangaTitle, sourceName, len(titles)),
		titles,
	)

	return selections, err
}
