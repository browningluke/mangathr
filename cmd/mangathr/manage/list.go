package manage

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/config"
	"github.com/browningluke/mangathr/v2/internal/database"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/utils"
	"strings"
)

func printList(driver *database.Driver, sourceFilter string, titleFilter []string) {
	allManga, err := driver.QueryAllManga()
	if err != nil {
		logging.ExitIfErrorWithFunc(&logging.ScraperError{
			Error: err, Message: "An error occurred while getting manga from database", Code: 0,
		}, closeDatabase)
	}

	for _, m := range allManga {
		if sourceFilter != "" {
			if strings.ToLower(sourceFilter) != strings.ToLower(m.Source) {
				continue
			}
		}

		if len(titleFilter) > 0 {
			if _, exists := utils.FindInSliceFold(titleFilter, m.Title); !exists {
				continue
			}
		}

		fmt.Printf(
			"\u001B[1m%s\u001B[0m\n"+
				"  Source:          %s\n"+
				"  Mapping:         %s\n"+
				"  Filtered Groups: %s\n"+
				"\n",
			m.Title,
			m.Source,
			m.Mapping,
			m.FilteredGroups,
		)
	}
}

func handleList(args *manageOpts, config *config.Config, driver *database.Driver) {
	printList(driver, args.List.Source, args.List.Query)
}
