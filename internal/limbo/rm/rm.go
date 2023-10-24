package rm

import (
	"fmt"
	"github.com/galqiwi/garc/internal/config"
	"github.com/galqiwi/garc/internal/limbo/common"
	"github.com/spf13/cobra"
)

var remoteUUID string

var LimboRmCmd = &cobra.Command{
	Use: "rm",
	RunE: func(cmd *cobra.Command, args []string) error {
		return limboRm()
	},
}

func init() {
	LimboRmCmd.PersistentFlags().StringVar(
		&remoteUUID,
		"uuid",
		"",
		"remote uuid",
	)
}

func limboRm() error {
	globalConfig, err := config.GetConfig()
	if err != nil {
		return err
	}
	limboClient := common.NewLimboClient(&globalConfig.LimboConfig)

	remoteExists, err := limboClient.DoesRemoteExist(remoteUUID)
	if err != nil {
		return err
	}
	if !remoteExists {
		fmt.Println("Remote does not exist, check with \"list\" command")
		return nil
	}

	remoteMeta, err := limboClient.GetRemoteMeta(remoteUUID)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Are you shure you want to delete %v? (y/n)\n", remoteMeta.Name)
	var answer string
	_, err = fmt.Scanln(&answer)
	if err != nil {
		return err
	}
	if answer != "yes" && answer != "y" {
		fmt.Println("exiting")
		return nil
	}

	err = limboClient.RemoveRemote(remoteUUID)
	if err != nil {
		return err
	}

	fmt.Printf("removed %v\n", remoteMeta.Name)

	return nil
}
