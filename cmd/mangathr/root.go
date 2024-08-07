package main

import (
	cmdConfig "github.com/browningluke/mangathr/v2/cmd/mangathr/config"
	"github.com/browningluke/mangathr/v2/cmd/mangathr/download"
	"github.com/browningluke/mangathr/v2/cmd/mangathr/manage"
	"github.com/browningluke/mangathr/v2/cmd/mangathr/register"
	"github.com/browningluke/mangathr/v2/cmd/mangathr/update"
	"github.com/browningluke/mangathr/v2/cmd/mangathr/version"
	"github.com/browningluke/mangathr/v2/internal/config"
	"github.com/browningluke/mangathr/v2/internal/config/defaults"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/browningluke/mangathr/v2/internal/ui"
	"github.com/browningluke/mangathr/v2/internal/utils"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfg      *config.Config
	cfgFile  string
	logLevel string

	rootCmd = &cobra.Command{
		Use:                   "mangathr [OPTIONS]",
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
		"", "Path to config file (default is $XDG_CONFIG_HOME/mangathr/config)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l",
		"", "Set the logging level (\"debug\"|\"info\"|\"warn\"|\"error\"|\"off\") (default \"off\")")

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
	rootCmd.AddCommand(manage.NewCmd(cfg))
	rootCmd.AddCommand(version.NewCmd(cfg))
	rootCmd.AddCommand(cmdConfig.NewCmd(getConfigPath))
}

func getConfigPath() string {
	if utils.IsRunningInContainer() {
		return defaults.ConfigPathDocker()
	}

	if cfgFile != "" {
		// Use config file from the flag.
		return cfgFile
	} else {
		configDir := defaults.ConfigDir()
		err := os.MkdirAll(configDir, os.ModePerm)
		cobra.CheckErr(err)

		return defaults.ConfigPath()
	}
}

func initConfig() {
	configPath := getConfigPath()

	exists, err := cfg.Load(configPath, utils.IsRunningInContainer())
	if err != nil {
		if exists {
			// If config file exists, and contains invalid data, exit.
			// NOTE: an error will only be raised if the config file is
			// missing **required** values. It will try to repair
			// optional values on its own.
			ui.Fatalf("Config file contains invalid data.\nReason: %s\n", err.Error())
		} else {
			// If config file does not exist, all default values will
			// be used and only a warning will be raised.
			logging.Warningln("Config file missing (or contains invalid yaml), ignoring. Error: ", err)
		}
	}

	setLogLevel(logLevel, cfg.LogLevel)
	logging.Debugln("Using config path: ", configPath)
}

func setLogLevel(logLevelArg, logLevelConf string) {
	//fmt.Println("log level arg: ", logLevelArg)
	//fmt.Println("log level cfg: ", logLevelConf)

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
	case "error":
		loggingLevel = logging.ERROR
	case "warn":
		loggingLevel = logging.WARNING
	case "info":
		loggingLevel = logging.INFO
	case "debug":
		loggingLevel = logging.DEBUG
	}
	logging.SetLoggingLevel(loggingLevel)
}
