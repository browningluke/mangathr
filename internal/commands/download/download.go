package download

import (
	"mangathrV2/internal/config"
	"mangathrV2/internal/sources/scrapers"
)

func Run(args *Args, config *config.Config) {
	scraper := scrapers.NewScraper(args.Plugin)
	scraper.ListChapters()
	scraper.SelectChapters()
}
