package cmd

import (
	"fmt"
	"github.com/galqiwi/garc/internal/config"
	"github.com/galqiwi/garc/internal/utils/shell"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var EditConfigCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit config",
	Long:  "Edits config, respects $EDITOR variable",
	RunE: func(cmd *cobra.Command, args []string) error {
		return editConfig()
	},
}

func editConfig() error {
	path, err := config.GetConfigPath()
	if err != nil {
		return err
	}

	editor, exists := os.LookupEnv("EDITOR")
	if !exists {
		if ok, err := shell.CommandExists("nano"); err == nil && ok {
			editor = "nano"
			exists = true
		}
	}

	if !exists {
		fmt.Println("$EDITOR is empty. Set it up like this:\nexport EDITOR=nano")
		return nil
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("%v %v", editor, path))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
