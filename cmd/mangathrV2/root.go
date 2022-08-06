package main

import (
	"github.com/browningluke/mangathrV2/cmd/mangathrV2/download"
	"github.com/browningluke/mangathrV2/cmd/mangathrV2/register"
	"github.com/browningluke/mangathrV2/cmd/mangathrV2/update"
	"github.com/browningluke/mangathrV2/internal/config"
	"github.com/browningluke/mangathrV2/internal/config/defaults"
	"github.com/browningluke/mangathrV2/internal/logging"
	"github.com/browningluke/mangathrV2/internal/utils"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfg      *config.Config
	cfgFile  string
	logLevel string

	rootCmd = &cobra.Command{
		Use:                   "mangathrv2 [OPTIONS]",
		Short:                 "📦 A CLI utility for downloading Manga & metadata.",
		DisableFlagsInUseLine: true,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Config
	cfg = &config.Config{}
	cobra.OnInitialize(initConfig)

	// Flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config",
		"", "Path to config file (default is $XDG_CONFIG_HOME/mangathrv2/config)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l",
		"off", "Set the logging level (\"debug\"|\"info\"|\"warn\"|\"error\"|\"off\")")

	// Help func
	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetHelpFunc(func(c *cobra.Command, s []string) {
		err := c.Usage()
		if err != nil {
			panic(err)
		}
	})
	rootCmd.SetHelpCommand(&cobra.Command{Use: "h", Hidden: true})
	rootCmd.PersistentFlags().BoolP("help", "h",
		false, "Print this text")

	// Sub commands
	rootCmd.AddCommand(download.NewCmd(cfg))
	rootCmd.AddCommand(register.NewCmd(cfg))
	rootCmd.AddCommand(update.NewCmd(cfg))
}

func initConfig() {
	filePath := ""
	if cfgFile != "" {
		// Use config file from the flag.
		filePath = cfgFile
	} else {
		configDir := defaults.ConfigDir()
		err := os.MkdirAll(configDir, os.ModePerm)
		cobra.CheckErr(err)

		filePath = defaults.ConfigPath()
	}

	err := cfg.Load(filePath, utils.IsRunningInContainer())
	cobra.CheckErr(err)

	setLogLevel(logLevel, cfg.LogLevel)

}

func setLogLevel(logLevelArg, logLevelConf string) {
	logging.Infoln("log level arg: ", logLevelArg)
	logging.Infoln("log level cfg: ", logLevelConf)

	// If neither value is set, do nothing (level has default: logging.loggingLevel)
	if logLevelArg == "" && logLevelConf == "" {
		return
	}

	selectedLevel := logLevelArg
	if selectedLevel == "" {
		// Use config (or default) as second priority
		selectedLevel = logLevelConf
	}

	loggingLevel := logging.OFF
	switch selectedLevel {
	case "ERROR":
		loggingLevel = logging.ERROR
	case "WARNING":
		loggingLevel = logging.WARNING
	case "INFO":
		loggingLevel = logging.INFO
	case "DEBUG":
		loggingLevel = logging.DEBUG
	}
	logging.SetLoggingLevel(loggingLevel)
}
