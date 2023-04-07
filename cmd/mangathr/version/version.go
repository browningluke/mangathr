package version

import (
	"fmt"
	"github.com/browningluke/mangathr/internal/config"
	"github.com/spf13/cobra"
)

func NewCmd(_ *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the version number of mangathr",
		Aliases: []string{"v"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("mangathr v2.0.0 -- HEAD")
		},
		DisableFlagsInUseLine: true,
	}

	return cmd
}
