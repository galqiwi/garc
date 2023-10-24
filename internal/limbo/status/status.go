package status

import (
	"fmt"
	"github.com/galqiwi/garc/internal/config"
	"github.com/galqiwi/garc/internal/limbo/common"
	"github.com/spf13/cobra"
)

var LimboStatusCmd = &cobra.Command{
	Use: "status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return limboStatus()
	},
}

func limboStatus() error {
	globalConfig, err := config.GetConfig()
	if err != nil {
		return err
	}
	limboClient := common.NewLimboClient(&globalConfig.LimboConfig)

	limboMeta, err := common.ReadCurrentArchiveMeta()
	if err != nil {
		return err
	}
	remoteUUID := fmt.Sprint(limboMeta.Id)

	fmt.Printf("%v\n\n", limboMeta.Name)
	fmt.Printf("current version is: %v\n", limboMeta.Version)

	remoteExists, err := limboClient.DoesRemoteExist(remoteUUID)
	if err != nil {
		return err
	}
	if !remoteExists {
		fmt.Print("remote does not exist\n\n")
		return nil
	}

	remoteMeta, err := limboClient.GetRemoteMeta(remoteUUID)
	if err != nil {
		return err
	}
	fmt.Printf("remote version is: %v\n\n", remoteMeta.Version)
	return nil
}
