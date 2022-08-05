package main

import (
	"fmt"
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
		Use:   "mangathrv2",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cfg = &config.Config{}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config",
		"", "config file (default is $XDG_CONFIG_HOME/mangathrv2/config)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l",
		"off", "Set the logging level (\"debug\"|\"info\"|\"warn\"|\"error\"|\"off\")")

	rootCmd.SetHelpFunc(func(c *cobra.Command, s []string) {
		err := c.Usage()
		if err != nil {
			panic(err)
		}
	})
	rootCmd.SetHelpCommand(&cobra.Command{Use: "h", Hidden: true})

	//rootCmd.AddCommand(versionCmd)
	//rootCmd.AddCommand(downloadCmd)

	//cobra.Command{
	//	Use:                        "",
	//	Aliases:                    nil,
	//	SuggestFor:                 nil,
	//	Short:                      "",
	//	Long:                       "",
	//	Example:                    "",
	//	ValidArgs:                  nil,
	//	ValidArgsFunction:          nil,
	//	Args:                       nil,
	//	ArgAliases:                 nil,
	//	BashCompletionFunction:     "",
	//	Deprecated:                 "",
	//	Annotations:                nil,
	//	Version:                    "",
	//	PersistentPreRun:           nil,
	//	PersistentPreRunE:          nil,
	//	PreRun:                     nil,
	//	PreRunE:                    nil,
	//	Run:                        nil,
	//	RunE:                       nil,
	//	PostRun:                    nil,
	//	PostRunE:                   nil,
	//	PersistentPostRun:          nil,
	//	PersistentPostRunE:         nil,
	//	FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	//	CompletionOptions:          cobra.CompletionOptions{},
	//	TraverseChildren:           false,
	//	Hidden:                     false,
	//	SilenceErrors:              false,
	//	SilenceUsage:               false,
	//	DisableFlagParsing:         false,
	//	DisableAutoGenTag:          false,
	//	DisableFlagsInUseLine:      false,
	//	DisableSuggestions:         false,
	//	SuggestionsMinimumDistance: 0,
	//}

	//rootCmd.SetHelpCommand(nil)

	//rootCmd.AddCommand(initCmd)
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
	fmt.Println("log level arg: ", logLevelArg)

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
