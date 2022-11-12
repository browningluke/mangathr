package manage

import (
	"errors"
	"fmt"
	"github.com/browningluke/mangathrV2/ent"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/database"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/sources"
	"github.com/browningluke/mangathrV2/internal/ui"
	"github.com/browningluke/mangathrV2/internal/utils"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *config.Config) *cobra.Command {
	o := &manageOpts{}

	cmd := &cobra.Command{
		Use:     "manage",
		Short:   "Manage series registered in database",
		Aliases: []string{"m"},
		Run: func(cmd *cobra.Command, args []string) {
			o.runWrapper(cfg, handleMenu)
		},
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(deleteSubcommand(cfg))

	return cmd
}

func deleteSubcommand(cfg *config.Config) *cobra.Command {
	o := &manageOpts{}

	cmd := &cobra.Command{
		Use:     "delete [OPTIONS] -s SOURCE QUERY [QUERY]...",
		Short:   "Delete series registered in database",
		Aliases: []string{"d"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires at least 1 series to remove")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Validate source flag input
			if _, exists := sources.MatchScraperTitle(o.Delete.Source); !exists {
				logging.ExitIfError(&logging.ScraperError{
					Error:   nil,
					Message: fmt.Sprintf("%s is not a valid source", o.Delete.Source), Code: 0,
				})
			}

			o.Delete.Query = args
			o.runWrapper(cfg, handleDelete)
		},
		DisableFlagsInUseLine: true,
	}

	cmd.Flags().StringVarP(&o.Delete.Source, "source", "s",
		"", "Source for desired series")
	err := cmd.MarkFlagRequired("source")
	cobra.CheckErr(err)

	return cmd
}

func (o *manageOpts) runWrapper(cfg *config.Config, f func(*manageOpts, *config.Config, *database.Driver)) {
	utils.CreateSigIntHandler(closeDatabase)

	// Open database
	var err error
	driver, err = database.GetDriver(database.SQLITE, cfg.Database.Uri)
	if err != nil {
		logging.Errorln(err)
		ui.Fatal("An error occurred while establishing a connection to the database")
	}
	defer closeDatabase()

	f(o, cfg, driver)
}
