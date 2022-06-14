package argparse

import (
	"github.com/akamensky/argparse"
	"mangathrV2/internal/commands/download"
	"mangathrV2/internal/commands/register"
	"os"
)

type Argparse struct {
	Command string

	// Commands
	Download download.Args
	Register register.Args

	Config struct {
	}

	// Options (overrides options in config file)
	Options struct {
		LogLevel string
		DryRun   bool
	}
}

func (a *Argparse) Parse() error {
	parser := argparse.NewParser("mangathr", "description")

	// Options
	loglevel := parser.Selector("l", "loglevel", []string{"INFO", "DEBUG", "WARN", "ERROR"},
		&argparse.Options{Help: "", Required: false})

	dryrun := parser.Flag("d", "dry-run", &argparse.Options{Help: "", Required: false})

	// Commands

	downloadCmd := parser.NewCommand("download", "")
	downloadPlugin := downloadCmd.Selector("p", "plugin", []string{"mangadex", "cubari"},
		&argparse.Options{Help: "", Required: true})
	downloadQuery := downloadCmd.String("q", "query", &argparse.Options{Help: "", Required: true})
	downloadAll := downloadCmd.Flag("a", "all", &argparse.Options{Help: "", Required: false})

	registerCmd := parser.NewCommand("register", "")
	registerPlugin := registerCmd.Selector("p", "plugin", []string{"mangadex", "cubari"},
		&argparse.Options{Help: "", Required: true})
	registerQuery := registerCmd.String("q", "query", &argparse.Options{Help: "", Required: true})
	registerYes := registerCmd.Flag("y", "yes", &argparse.Options{Help: "", Required: false})

	updateCmd := parser.NewCommand("update", "")

	err := parser.Parse(os.Args)
	if err != nil {
		return err
	}

	a.Options.LogLevel = *loglevel
	a.Options.DryRun = *dryrun
	if downloadCmd.Happened() {
		a.Command = "download"
		a.Download.Plugin = *downloadPlugin
		a.Download.Query = *downloadQuery
		a.Download.All = *downloadAll

	} else if registerCmd.Happened() {
		a.Command = "register"
		a.Register.Query = *registerQuery
		a.Register.Yes = *registerYes

		a.Register.Plugin = *registerPlugin
	} else if updateCmd.Happened() {
		a.Command = "update"
	}

	return nil
}
