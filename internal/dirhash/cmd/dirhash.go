package cmd

import (
	"github.com/galqiwi/garc/internal/dirhash/ls"
	"github.com/spf13/cobra"
)

var DirHashCmd = &cobra.Command{
	Use:   "dirhash",
	Short: "utils for recursive directory hashing and verification",
}

func init() {
	DirHashCmd.AddCommand(ls.LsCmd)
}
