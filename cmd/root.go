package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("spinNamespace", "")
}

func GetRootCmd(args []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "toolctl",
		Short:             "tool control interface.",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}

	rootCmd.SetArgs(args)

	rootCmd.AddCommand(spinCmd())

	return rootCmd
}
