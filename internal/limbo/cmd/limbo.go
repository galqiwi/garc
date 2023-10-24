package cmd

import (
	"github.com/galqiwi/garc/internal/limbo/limbo_init"
	"github.com/galqiwi/garc/internal/limbo/list"
	"github.com/galqiwi/garc/internal/limbo/pull"
	"github.com/galqiwi/garc/internal/limbo/push"
	"github.com/galqiwi/garc/internal/limbo/rm"
	"github.com/galqiwi/garc/internal/limbo/status"
	"github.com/spf13/cobra"
)

var LimboCmd = &cobra.Command{
	Use:   "limbo",
	Short: "Utils for moving non-etomb files to limbo",
}

func init() {
	LimboCmd.AddCommand(limbo_init.LimboInitCmd)
	LimboCmd.AddCommand(push.LimboPushCmd)
	LimboCmd.AddCommand(list.LimboListCmd)
	LimboCmd.AddCommand(pull.LimboPullCmd)
	LimboCmd.AddCommand(rm.LimboRmCmd)
	LimboCmd.AddCommand(status.LimboStatusCmd)
}
