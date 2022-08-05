package update

import (
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/database"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/ui"
	"github.com/spf13/cobra"
	"time"
)

func NewCmd(cfg *config.Config) *cobra.Command {
	o := &updateOpts{}

	cmd := &cobra.Command{
		Use:   "update [OPTIONS]",
		Short: "Check for new chapters",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.Downloader.DryRun = cfg.Downloader.DryRun || o.DryRun
			o.run(cfg)
		},
		DisableFlagsInUseLine: true,
	}

	cmd.Flags().BoolVarP(&o.DryRun, "dry-run", "",
		false, "do not download files or update database")

	return cmd
}

func (o *updateOpts) run(config *config.Config) {
	// Open database
	var err error
	driver, err = database.GetDriver(database.SQLITE, config.Database.Uri)
	if err != nil {
		logging.Errorln(err)
		ui.Fatal("Unable to open database.")
	}
	defer closeDatabase()

	allManga, err := driver.QueryAllManga()
	if err != nil {
		logging.ExitIfErrorWithFunc(&logging.ScraperError{
			Error: err, Message: "An error occurred while getting Manga from database.", Code: 0,
		}, closeDatabase)
	}

	for _, manga := range allManga {
		checkMangaForNewChapters(config, manga)

		// Sleep between checks
		dur, durErr := time.ParseDuration(config.Downloader.Delay.UpdateChapter)
		if durErr != nil {
			logging.Errorln(durErr)
			ui.Error("Failed to parse time duration.")
		}
		time.Sleep(dur)
	}
}
