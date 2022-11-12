package manage

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
	o := &manageOpts{}

	cmd := &cobra.Command{
		Use:     "manage",
		Short:   "Manage series registered in database",
		Aliases: []string{"m"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New("manage does not accept arguments")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.run(cfg)
		},
		DisableFlagsInUseLine: true,
	}

	return cmd
}

func (o *manageOpts) run(cfg *config.Config) {
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
func deleteFromDatabase(filter func(manga *ent.Manga) bool) {
	allManga, err := driver.QueryAllManga()
	if err != nil {
		logging.ExitIfErrorWithFunc(&logging.ScraperError{
			Error: err, Message: "An error occurred while getting manga from database", Code: 0,
		}, closeDatabase)
	}

	for _, manga := range allManga {
		if filter(manga) {
			err := driver.DeleteManga(manga)
			if err != nil {
				panic(err)
			}
		}
	}
}
