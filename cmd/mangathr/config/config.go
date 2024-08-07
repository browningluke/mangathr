package config

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/config"
	"github.com/browningluke/mangathr/v2/internal/ui"
	"github.com/browningluke/mangathr/v2/internal/utils"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func NewCmd(getConfigPath func() string) *cobra.Command {
	o := &configOpts{}

	cmd := &cobra.Command{
		Use:     "config [-e|--edit] | [-s|--show] | [-p|--path] ",
		Short:   "Manage the config file",
		Aliases: []string{"c"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			configPath := getConfigPath()

			// -e
			if o.Edit {
				err := editConfig(configPath)
				if err != nil {
					panic(err)
				}
			}

			// -s
			if o.Show {
				err := showFile(configPath)
				if err != nil {
					panic(err)
				}
			}

			// -p
			if o.Path {
				fmt.Println(configPath)
			}
		},
		DisableFlagsInUseLine: true,
	}

	cmd.Flags().BoolVarP(&o.Edit, "edit", "e", false, "Opens an editor to modify the config file")
	cmd.Flags().BoolVarP(&o.Show, "show", "s", false, "Prints the config file")
	cmd.Flags().BoolVarP(&o.Path, "path", "p", false, "Shows the path to the config file")
	cmd.MarkFlagsMutuallyExclusive("edit", "show", "path")

	return cmd
}

func editConfig(filePath string) error {
	var cmd *exec.Cmd

	// Determine the default editor based on the operating system
	switch goos := runtime.GOOS; goos {
	case "windows":
		cmd = exec.Command("notepad.exe", filePath)
	case "darwin":
		fallthrough
	case "linux":
		fallthrough
	default:
		cmd = exec.Command("vi", filePath)
	}

	// Sourced from: https://stackoverflow.com/a/12089980
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Open editor
	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	// Validate config
	cfg := config.Config{}
	exists, err := cfg.Load(filePath, utils.IsRunningInContainer())
	if err != nil {
		// if file doesn't exist, something else went wrong
		if !exists {
			return err
		}

		// Config validation failed, force user to re-edit
		ui.Errorf("Config invalid: %s", err)
		fmt.Println("\nPress Enter to re-edit...")
		_, err := fmt.Scanln()
		if err != nil {
			return err
		}

		return editConfig(filePath)
	}

	ui.PrintlnColor(ui.Green, "Config valid, saving...")
	return nil
}

func showFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return err
	}

	return nil
}
