package cmd

import "github.com/spf13/cobra"

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Config utils",
}

func init() {
	ConfigCmd.AddCommand(PrintConfigCmd)
	ConfigCmd.AddCommand(InitConfigCmd)
	ConfigCmd.AddCommand(UpdateConfigCmd)
	ConfigCmd.AddCommand(EditConfigCmd)
}
