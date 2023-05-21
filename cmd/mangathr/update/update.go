package update

import (
	"fmt"
	"github.com/browningluke/mangathr/internal/config"
	"github.com/browningluke/mangathr/internal/database"
	"github.com/browningluke/mangathr/internal/logging"
	"github.com/browningluke/mangathr/internal/ui"
	"github.com/browningluke/mangathr/internal/utils"
	"github.com/spf13/cobra"
	"time"
)

type seriesStats struct {
	found      int
	downloaded int
	errors     int
}

type updateStats struct {
	checked       int
	foundChapters int
	foundSeries   int
	downloaded    int
	errors        int
}

func NewCmd(cfg *config.Config) *cobra.Command {
	o := &updateOpts{}

	cmd := &cobra.Command{
		Use:     "update [OPTIONS]",
		Short:   "Check for new chapters",
		Aliases: []string{"u"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.Downloader.DryRun = cfg.Downloader.DryRun || o.DryRun
			o.run(cfg)
		},
		DisableFlagsInUseLine: true,
	}

	cmd.Flags().BoolVarP(&o.DryRun, "dry-run", "",
		false, "Do not download files or update database")

	return cmd
}

func printStats(stats updateStats) {
	fmt.Printf(
		"\n\033[1mChecked:    \033[0m%d\n"+
			"\u001B[1mFound:      \u001B[0m%d for %d series\n"+
			"\u001B[1mDownloaded: \u001B[0m%d/%d\n"+
			"\u001B[1mErrors:     \u001B[0m%d\n",
		stats.checked,
		stats.foundChapters, stats.foundSeries,
		stats.downloaded, stats.foundChapters,
		stats.errors,
	)
}

func (o *updateOpts) run(cfg *config.Config) {
	// Propagate config to all sub-configs
	cfg.Propagate()

	utils.CreateSigIntHandler(closeDatabase)

	// Open database
	var err error
	driver, err = database.GetDriver(database.SQLITE, cfg.Database.Uri)
	if err != nil {
		logging.Errorln(err)
		ui.Fatalf("Unable to open database.\nReason: %s\n", err)
	}
	defer closeDatabase()

	allManga, err := driver.QueryAllManga()
	if err != nil {
		logging.ExitIfErrorWithFunc(&logging.ScraperError{
			Error: err, Message: "An error occurred while loading all Manga from database.", Code: 0,
		}, closeDatabase)
	}

	stats := updateStats{}

	for _, manga := range allManga {
		s := checkMangaForNewChapters(manga)

		// Update stats
		stats.checked++
		if s.found > 0 {
			stats.downloaded += s.downloaded
			stats.foundSeries++
			stats.foundChapters += s.found
			stats.errors += s.errors
		}

		// Sleep between checks
		dur, durErr := time.ParseDuration(cfg.Downloader.Delay.UpdateChapter)
		if durErr != nil {
			logging.Errorln(durErr)
			ui.Error("Failed to parse time duration.")
		}
		time.Sleep(dur)
	}

	printStats(stats)
}
