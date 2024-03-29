package manage

import (
	"errors"
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/config"
	"github.com/browningluke/mangathr/v2/internal/database"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/sources"
	"github.com/browningluke/mangathr/v2/internal/ui"
	"github.com/browningluke/mangathr/v2/internal/utils"
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
	cmd.AddCommand(listSubcommand(cfg))

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

func listSubcommand(cfg *config.Config) *cobra.Command {
	o := &manageOpts{}

	cmd := &cobra.Command{
		Use:     "list [-s SOURCE] [QUERY]...",
		Short:   "List series registered in database",
		Aliases: []string{"l"},
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if o.List.Source != "" {
				// Validate source flag input
				if _, exists := sources.MatchScraperTitle(o.Delete.Source); !exists {
					logging.ExitIfError(&logging.ScraperError{
						Error:   nil,
						Message: fmt.Sprintf("%s is not a valid source", o.Delete.Source), Code: 0,
					})
				}
			}

			o.List.Query = args
			o.runWrapper(cfg, handleList)
		},
		DisableFlagsInUseLine: true,
	}

	cmd.Flags().StringVarP(&o.Delete.Source, "source", "s",
		"", "Source for desired series")

	return cmd
}

func (o *manageOpts) runWrapper(cfg *config.Config, f func(*manageOpts, *config.Config, *database.Driver)) {
	// Propagate config to all sub-configs
	cfg.Propagate()

	utils.CreateSigIntHandler(closeDatabase)

	// Open database
	var err error
	driver, err = database.GetDriver()
	if err != nil {
		logging.Errorln(err)
		ui.Fatalf("Unable to open database.\nReason: %s\n", err)
	}
	defer closeDatabase()

	f(o, cfg, driver)
}
