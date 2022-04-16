package main

import (
	_ "github.com/mattn/go-sqlite3"
	"mangathrV2/internal/argparse"
	"mangathrV2/internal/commands/download"
	"mangathrV2/internal/commands/register"
	"mangathrV2/internal/config"
	"mangathrV2/internal/logging"
	"mangathrV2/internal/utils"
)

func main() {
	// Load config object, returns Config struct
	var c config.Config
	if err := c.Load("./examples/config.yml"); err != nil {
		utils.RaiseError(err)
	}

	// Load argparse object, returns ArgParse struct
	var a argparse.Argparse
	if err := a.Parse(); err != nil {
		utils.RaiseError(err)
	}

	// Init logging
	logging.Init()
	loggingLevel := logging.WARNING
	if a.Options.LogLevel != "" {
		switch a.Options.LogLevel {
		case "ERROR":
			loggingLevel = logging.ERROR
			break
		case "WARN":
			loggingLevel = logging.WARNING
			break
		case "INFO":
			loggingLevel = logging.INFO
			break
		case "DEBUG":
			loggingLevel = logging.DEBUG
			break
		}
	} else {
		switch c.LogLevel {
		case "ERROR":
			loggingLevel = logging.ERROR
			break
		case "WARN":
			loggingLevel = logging.WARNING
			break
		case "INFO":
			loggingLevel = logging.INFO
			break
		case "DEBUG":
			loggingLevel = logging.DEBUG
			break
		}
	}
	logging.SetLoggingLevel(loggingLevel)

	logging.Debugln(c)
	logging.Debugln(a)

	switch a.Command {
	case "download":
		logging.Infoln("Downloading", a.Download)
		download.Run(&a.Download, &c)
		break
	case "register":
		logging.Infoln("Registering", a.Register)
		register.Run(&a.Register, &c)
	}

	// Merge Config & ArgParse (ArgParse priority) into ProgramOptions
	// Call (download|register|update|manage|config).go > run(ProgramOptions po) to start program execution
}
