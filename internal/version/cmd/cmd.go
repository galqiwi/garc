package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime/debug"
	"time"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of this binary",
	RunE: func(cmd *cobra.Command, args []string) error {
		return versionCmd()
	},
}

func getCommit() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value
		}
	}
	return ""
}

func getTime() *time.Time {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.time" {
			output, err := time.Parse(time.RFC3339, setting.Value)
			if err != nil {
				return nil
			}
			return &output
		}
	}
	return nil
}

func versionCmd() error {
	commit := getCommit()
	commitTime := getTime()
	if commit == "" {
		fmt.Println("unknown")
		return nil
	}
	fmt.Printf("git commit %v\n", commit)
	if commitTime == nil {
		return nil
	}
	fmt.Printf("commit time is %v\n", commitTime.In(time.Now().Location()))
	return nil
}
