package register

import (
	"errors"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/database"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/ui"
	"github.com/browningluke/mangathrV2/internal/utils"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *config.Config) *cobra.Command {
	o := &registerOpts{}

	cmd := &cobra.Command{
		Use:     "register [OPTIONS] -s SOURCE QUERY",
		Short:   "Register chapters to database",
		Aliases: []string{"r"},
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

func (o *registerOpts) run(cfg *config.Config) {
	utils.CreateSigIntHandler(closeDatabase)

	// Open database
	var err error
	driver, err = database.GetDriver(database.SQLITE, cfg.Database.Uri)
	if err != nil {
		logging.Errorln(err)
		ui.Fatal("An error occurred while establishing a connection to the database")
	}
	defer closeDatabase()

	handleMenu(o, cfg, driver)
}
