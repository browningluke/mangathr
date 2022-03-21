package argparse

import (
	"github.com/akamensky/argparse"
	"os"
)

type Argparse struct {
	Command string
	Plugin  string
	Query   string
	All     bool
}

func (a *Argparse) Parse() error {
	parser := argparse.NewParser("mangathr", "description")

	downloadCmd := parser.NewCommand("download", "")
	downloadPlugin := downloadCmd.Selector("p", "plugin", []string{"mangadex", "webtoons"},
		&argparse.Options{Help: "", Required: true})
	downloadQuery := downloadCmd.String("q", "query", &argparse.Options{Help: "", Required: true})
	downloadAll := downloadCmd.Flag("a", "all", &argparse.Options{Help: "", Required: false})

	registerCmd := parser.NewCommand("register", "")
	registerPlugin := registerCmd.Selector("p", "plugin", []string{"mangadex", "webtoons"},
		&argparse.Options{Help: "", Required: true})
	registerQuery := registerCmd.String("q", "query", &argparse.Options{Help: "", Required: true})
	registerYes := registerCmd.Flag("y", "yes", &argparse.Options{Help: "", Required: false})

	err := parser.Parse(os.Args)
	if err != nil {
		return err
	}

	if downloadCmd.Happened() {
		a.Command = "download"
		a.Query = *downloadQuery
		a.All = *downloadAll

		a.Plugin = *downloadPlugin
	} else if registerCmd.Happened() {
		a.Command = "register"
		a.Query = *registerQuery
		a.All = *registerYes

		a.Plugin = *registerPlugin
	} else {
		a.Plugin = ""
	}

	return nil
}
