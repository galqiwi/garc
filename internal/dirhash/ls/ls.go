package ls

import (
	"errors"
	"fmt"
	"os"

	"github.com/galqiwi/garc/internal/dirhash/utils"
	"github.com/spf13/cobra"
)

var directory string
var output string

var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list directory contents with hashes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return lsCmd()
	},
}

func init() {
	LsCmd.Flags().StringVarP(
		&directory,
		"directory",
		"d",
		"",
		"directory to recursively hash",
	)
	LsCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"",
		"output file (if not specified, the output will be printed to stdout)",
	)
}

func lsCmd() error {
	if directory == "" {
		return errors.New("directory is required")
	}

	hashByPath, err := utils.GetHashMeta(directory)
	if err != nil {
		return fmt.Errorf("error getting hash meta: %w", err)
	}

	yamlOutput, err := utils.HashToYaml(hashByPath)
	if err != nil {
		return fmt.Errorf("error converting to YAML: %w", err)
	}

	if output == "" {
		fmt.Print(yamlOutput)
		return nil
	}

	err = os.WriteFile(output, []byte(yamlOutput), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
