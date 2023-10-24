package push

import (
	"errors"
	"fmt"
	"github.com/galqiwi/garc/internal/config"
	"github.com/galqiwi/garc/internal/limbo/common"
	io2 "github.com/galqiwi/garc/internal/utils/io"
	"github.com/galqiwi/garc/internal/utils/tarball"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"time"
)

var force bool

var LimboPushCmd = &cobra.Command{
	Use: "push",
	RunE: func(cmd *cobra.Command, args []string) error {
		return limboPush()
	},
}

func init() {
	LimboPushCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "force push")
}

func askAndRemoveRemote(client *common.LimboClient, remoteUUID string) error {
	if !force {
		remoteMeta, err := client.GetRemoteMeta(remoteUUID)
		if err != nil {
			return err
		}
		fmt.Printf(
			"Looks like remote %v already exists, are you shure you want to delete it? (y/n)\n",
			remoteMeta.Name,
		)
		var answer string
		_, err = fmt.Scanln(&answer)
		if err != nil {
			return err
		}
		if answer != "yes" && answer != "y" {
			fmt.Println("exiting")
			os.Exit(1)
			return nil
		}
	}

	err := client.RemoveRemote(remoteUUID)
	fmt.Println("deleted remote dir")
	return err
}

func limboPush() error {
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

	remoteExists, err := limboClient.DoesRemoteExist(remoteUUID)
	if err != nil {
		return err
	}
	if remoteExists {
		remoteMeta, err := limboClient.GetRemoteMeta(remoteUUID)
		if err != nil {
			return err
		}
		if remoteMeta.Version > limboMeta.Version {
			fmt.Println("Remote is ahead of local version. Please, compare them manually using pull.")
			fmt.Printf("remote version: %v\n", remoteMeta.Version)
			fmt.Printf("local version: %v\n", limboMeta.Version)
			return errors.New("remote is ahead of local")
		}

		err = askAndRemoveRemote(limboClient, remoteUUID)
		if err != nil {
			return err
		}
	}

	limboMeta.Version += 1
	limboMeta.ModificationTime = time.Now()
	err = common.WriteCurrentArchiveMeta(limboMeta)
	if err != nil {
		return err
	}

	toIgnore := common.GetLimboIgnore("./")

	tarballSizeCounter := &io2.CountingWriter{}
	err = tarball.CreateTarball("./", tarballSizeCounter, toIgnore)
	if err != nil {
		return err
	}
	tarballSize := tarballSizeCounter.BytesWritten()

	bar := progressbar.DefaultBytes(
		tarballSize,
		"uploading",
	)

	tarReader, tarWriter := io.Pipe()
	g := new(errgroup.Group)

	g.Go(func() error {
		defer func() {
			_ = tarWriter.Close()
		}()
		tarWriter := io.MultiWriter(tarWriter, bar)
		return tarball.CreateTarball("./", tarWriter, toIgnore)
	})

	g.Go(func() error {
		return limboClient.CreateRemote(limboMeta, tarReader)
	})

	return g.Wait()
}
