package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var VerifyCmd = &cobra.Command{
	Use:           "verify",
	Short:         "verify that system is clean",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return verify()
	},
}

func verify() error {
	messages, err := getMessages()
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	displayMessages(messages)

	return fmt.Errorf("")
}
