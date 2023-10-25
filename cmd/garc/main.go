package main

import (
	"fmt"
	cmd6 "github.com/galqiwi/garc/internal/compess/cmd"
	"github.com/galqiwi/garc/internal/config"
	cmd2 "github.com/galqiwi/garc/internal/config/cmd"
	"github.com/galqiwi/garc/internal/limbo/cmd"
	cmd5 "github.com/galqiwi/garc/internal/ls/cmd"
	cmd3 "github.com/galqiwi/garc/internal/update/cmd"
	cmd4 "github.com/galqiwi/garc/internal/version/cmd"
	"os"

	"github.com/spf13/cobra"
)

var GarcCmd = &cobra.Command{
	Use:   "garc",
	Short: "Utils for archive management",
}

func init() {
	GarcCmd.AddCommand(cmd.LimboCmd)
	GarcCmd.AddCommand(cmd2.ConfigCmd)
	GarcCmd.AddCommand(cmd3.UpdateCmd)
	GarcCmd.AddCommand(cmd4.VersionCmd)
	GarcCmd.AddCommand(cmd5.LsCmd)
	GarcCmd.AddCommand(cmd6.CompressCmd)
	config.AddConfigFlag(GarcCmd)
}

func main() {
	err := Main()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Main() error {
	return GarcCmd.Execute()
}
