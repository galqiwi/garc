package cmd

import (
	"github.com/galqiwi/garc/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

var PrintConfigCmd = &cobra.Command{
	Use:   "print",
	Short: "Print config",
	RunE: func(cmd *cobra.Command, args []string) error {
		return printConfig()
	},
}

func printConfig() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}
	return yaml.NewEncoder(os.Stdout).Encode(cfg)
}
