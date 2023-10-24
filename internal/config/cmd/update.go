package cmd

import (
	"fmt"
	"github.com/galqiwi/garc/internal/config"
	"github.com/spf13/cobra"
)

var UpdateConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "Update config to latest version (useful for development)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateConfig()
	},
}

func updateConfig() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}
	err = config.SaveConfig(cfg)
	if err != nil {
		return err
	}
	path, err := config.GetConfigPath()
	if err != nil {
		return err
	}
	fmt.Printf("Updated config at %v\n", path)
	return nil
}
