package pull

import (
	"fmt"
	"github.com/galqiwi/garc/internal/config"
	"github.com/galqiwi/garc/internal/limbo/common"
	"github.com/galqiwi/garc/internal/utils/tarball"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"io"
)

var remoteUUID string
var destinationDir string

var LimboPullCmd = &cobra.Command{
	Use: "pull",
	RunE: func(cmd *cobra.Command, args []string) error {
		return limboPull()
	},
}

func init() {
	LimboPullCmd.PersistentFlags().StringVar(
		&remoteUUID,
		"uuid",
		"",
		"remote uuid",
	)
	LimboPullCmd.PersistentFlags().StringVar(
		&destinationDir,
		"dir",
		"",
		"destination directory",
	)
}

func limboPull() error {
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

	tarballSize, err := limboClient.GetRemoteSize(remoteUUID)
	if err != nil {
		return err
	}

	tarReader, tarWriter := io.Pipe()

	bar := progressbar.DefaultBytes(
		tarballSize,
		"downloading",
	)

	g := new(errgroup.Group)

	g.Go(func() error {
		defer func() {
			_ = tarWriter.Close()
		}()
		return limboClient.ReadRemoteTarball(
			remoteUUID,
			io.MultiWriter(bar, tarWriter),
		)
	})

	if destinationDir == "" {
		meta, err := limboClient.GetRemoteMeta(remoteUUID)
		if err != nil {
			return err
		}
		destinationDir = meta.Name
	}

	g.Go(func() error {
		return tarball.ExtractTarball(tarReader, destinationDir)
	})

	return g.Wait()
}
