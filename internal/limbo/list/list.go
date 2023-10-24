package list

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/galqiwi/garc/internal/config"
	"github.com/galqiwi/garc/internal/limbo/common"
	"github.com/google/uuid"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"math"
	"sort"
	"time"
)

var LimboListCmd = &cobra.Command{
	Use: "list",
	RunE: func(cmd *cobra.Command, args []string) error {
		return limboList()
	},
}

func prettyByteSize(b int64) string {
	bf := float64(b)
	for _, unit := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
		if math.Abs(bf) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", bf, unit)
		}
		bf /= 1024.0
	}
	return fmt.Sprintf("%.1fYiB", bf)
}

func limboList() error {
	globalConfig, err := config.GetConfig()
	if err != nil {
		return err
	}
	limboClient := common.NewLimboClient(&globalConfig.LimboConfig)

	remoteUUIDs, err := limboClient.ListRemotes()
	if err != nil {
		return err
	}

	var remoteMetas []*common.ArchiveMeta
	remoteSizeByUUID := make(map[uuid.UUID]int64)
	for _, remoteUUID := range remoteUUIDs {
		meta, err := limboClient.GetRemoteMeta(remoteUUID)
		if err != nil {
			return err
		}
		remoteMetas = append(remoteMetas, meta)

		tarSize, err := limboClient.GetRemoteSize(remoteUUID)
		if err != nil {
			return err
		}
		remoteSizeByUUID[meta.Id] = tarSize
	}

	sort.Slice(remoteMetas, func(i, j int) bool {
		return remoteMetas[j].ModificationTime.Before(remoteMetas[i].ModificationTime)
	})

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("UUID", "Name", "Size", "Age", "Modified", "Version")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, meta := range remoteMetas {
		tarSize := remoteSizeByUUID[meta.Id]
		tbl.AddRow(
			meta.Id,
			meta.Name,
			prettyByteSize(tarSize),
			time.Since(meta.CreationTime).Round(time.Second),
			fmt.Sprintf("%v ago", time.Since(meta.ModificationTime).Round(time.Second)),
			meta.Version,
		)
	}

	tbl.Print()

	return nil
}
