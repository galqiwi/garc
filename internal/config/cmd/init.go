package cmd

import (
	"fmt"
	"github.com/galqiwi/garc/internal/config"
	"github.com/spf13/cobra"
)

var InitConfigCmd = &cobra.Command{
	Use:   "limbo_init",
	Short: "Initialize empty config",
	RunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
}

func initConfig() error {
	err := config.SaveConfig(&config.Config{})
	if err != nil {
		return err
	}
	path, err := config.GetConfigPath()
	if err != nil {
		return err
	}
	fmt.Printf("Created empty config at %v\n", path)
	return nil
}
