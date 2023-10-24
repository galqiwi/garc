package limbo_init

import (
	"fmt"
	"github.com/galqiwi/garc/internal/limbo/common"
	"github.com/spf13/cobra"
)

var archiveName string

var LimboInitCmd = &cobra.Command{
	Use: "init",
	RunE: func(cmd *cobra.Command, args []string) error {
		return limboInit()
	},
}

func init() {
	LimboInitCmd.PersistentFlags().StringVar(
		&archiveName,
		"name",
		"",
		"archive name",
	)
}

func limboInit() error {
	if _, err := common.ReadCurrentArchiveMeta(); err == nil {
		return fmt.Errorf("is already a limbo directory")
	}

	err := common.WriteCurrentArchiveMeta(common.NewArchiveMeta(archiveName))
	if err != nil {
		return err
	}

	return nil
}
