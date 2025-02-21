package diff

import (
	"fmt"
	"os"

	"github.com/galqiwi/garc/internal/dirhash/utils"
	"github.com/spf13/cobra"
)

var directory string
var dirHashFile string

var DiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "diff between real dirhash and dumped dirhash",
	RunE: func(cmd *cobra.Command, args []string) error {
		return diffCmd()
	},
}

func init() {
	DiffCmd.Flags().StringVarP(&directory, "directory", "d", "", "directory to hash")
	DiffCmd.Flags().StringVarP(&dirHashFile, "dirhash-file", "", "", "file containing the dirhash")
}

func diffCmd() error {
	dirHashBytes, err := os.ReadFile(dirHashFile)
	if err != nil {
		return err
	}

	dirHash, err := utils.YamlToHash(string(dirHashBytes))
	if err != nil {
		return err
	}

	realDirHash, err := utils.GetHashMeta(directory)
	if err != nil {
		return err
	}

	delta := utils.GetDelta(dirHash, realDirHash)

	fmt.Println("Deleted:")
	for path, hash := range delta.Deleted {
		fmt.Printf("  %s: %s\n", path, hash)
	}

	fmt.Println("New:")
	for path, hash := range delta.New {
		fmt.Printf("  %s: %s\n", path, hash)
	}

	fmt.Println("ChangedOld:")
	for path, hash := range delta.ChangedOld {
		fmt.Printf("  %s: %s\n", path, hash)
	}

	fmt.Println("ChangedNew:")
	for path, hash := range delta.ChangedNew {
		fmt.Printf("  %s: %s\n", path, hash)
	}

	return nil
}
