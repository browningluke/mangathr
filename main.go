package main

import (
	"github.com/browningluke/mangathrV2/internal/argparse"
	"github.com/browningluke/mangathrV2/internal/commands/download"
	"github.com/browningluke/mangathrV2/internal/commands/register"
	"github.com/browningluke/mangathrV2/internal/commands/update"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/config/defaults"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/ui"
	"os"
)

/*
	Config
*/
func getConfigPath() (string, error) {
	// Create config directory if it does not exist
	configDir := defaults.ConfigDir()
	err := os.MkdirAll(configDir, os.ModePerm)
	return defaults.ConfigPath(), err
}

func loadConfig() (config.Config, error) {
	// Load config object, returns Config struct
	path, err := getConfigPath()
	if err != nil {
		return config.Config{}, err
	}

	var c config.Config
	err = c.Load(path)
	return c, err
}

/*
	Argparse
*/

func parseArgs() (argparse.Argparse, error) {
	// Load argparse object, returns ArgParse struct
	var a argparse.Argparse
	err := a.Parse()
	return a, err
}

func setLogLevel(logLevelArg, logLevelConf string) {
	// If neither value is set, do nothing (level has default: logging.loggingLevel)
	if logLevelArg == "" && logLevelConf == "" {
		return
	}

	selectedLevel := logLevelArg
	if selectedLevel == "" {
		// Use config (or default) as second priority
		selectedLevel = logLevelConf
	}

	var loggingLevel logging.Level
	switch selectedLevel {
	case "ERROR":
		loggingLevel = logging.ERROR
		break
	case "WARNING":
		loggingLevel = logging.WARNING
		break
	case "INFO":
		loggingLevel = logging.INFO
		break
	case "DEBUG":
		loggingLevel = logging.DEBUG
		break
	}
	logging.SetLoggingLevel(loggingLevel)
}

func main() {
	/*
		Parse config
	*/
	c, err := loadConfig()

	// If an error, we revert to using defaults
	if err != nil {
		logging.Warningln(err)
	}

	/*
		Parse args
	*/
	a, err := parseArgs()
	if err != nil {
		ui.Fatalf("%s%s\n", "Invalid arguments: ", err)
	}

	// Set log level
	setLogLevel(a.Options.LogLevel, c.LogLevel)

	// Prioritize dryrun arg over config setting
	c.Downloader.DryRun = a.Options.DryRun

	switch a.Command {
	case "download":
		logging.Infoln("Downloading", a.Download)
		download.Run(&a.Download, &c)
		break
	case "register":
		logging.Infoln("Registering", a.Register)
		register.Run(&a.Register, &c)
	case "update":
		logging.Infoln("Updating")
		update.Run(&c)
	}
}
