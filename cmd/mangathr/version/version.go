package version

import (
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/config"
	"github.com/browningluke/mangathr/v2/internal/utils"
	"github.com/browningluke/mangathr/v2/internal/version"
	"github.com/spf13/cobra"
)

func NewCmd(_ *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the version number of mangathr",
		Aliases: []string{"v"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			containerStr := ""

			if utils.IsRunningInContainer() {
				containerStr = "-docker"
			}

			fmt.Printf(
				"mangathr %s%s -- %s\n",
				version.GetVersion(),
				containerStr,
				version.GetSHA(),
			)
		},
		DisableFlagsInUseLine: true,
	}

	return cmd
}
