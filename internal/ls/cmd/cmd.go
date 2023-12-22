package cmd

import (
	"github.com/fatih/color"
	"github.com/galqiwi/garc/internal/utils/misc"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

var displayLongestPath bool
var sortBySize bool
var sortByPathLength bool
var sortByNFiles bool

var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "advanced directory listing",
	RunE: func(cmd *cobra.Command, args []string) error {
		return lsCmd()
	},
}

type entryMeta struct {
	Name           string
	Mode           os.FileMode
	MaxPathLength  int
	MaxPathExample string
	Size           int64
	NFiles         int64
}

func init() {
	LsCmd.PersistentFlags().BoolVarP(
		&displayLongestPath,
		"display-longest-path",
		"",
		false,
		"display longest path",
	)
	LsCmd.PersistentFlags().BoolVarP(
		&sortBySize,
		"size",
		"s",
		false,
		"sort by size",
	)
	LsCmd.PersistentFlags().BoolVarP(
		&sortByNFiles,
		"n-files",
		"n",
		false,
		"sort by # of files",
	)
	LsCmd.PersistentFlags().BoolVarP(
		&sortByPathLength,
		"path-length",
		"l",
		false,
		"sort by path length",
	)
}

func getFileMeta(root string, e os.DirEntry) (*entryMeta, error) {
	path := filepath.Join(root, e.Name())
	maxPathLength := len(path)

	info, err := e.Info()
	if err != nil {
		return nil, err
	}

	return &entryMeta{
		Name:           info.Name(),
		Mode:           info.Mode(),
		MaxPathLength:  maxPathLength,
		Size:           info.Size(),
		MaxPathExample: path,
		NFiles:         1,
	}, nil
}

func getEntryMeta(root string, e os.DirEntry) (*entryMeta, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	if !e.IsDir() {
		return getFileMeta(root, e)
	}

	dirPath := filepath.Join(root, e.Name())

	maxPathLength := len(dirPath)
	maxPathExample := dirPath
	var size int64 = 0
	var nFiles int64 = 0

	err = filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		nFiles += 1

		if maxPathLength < len(path) {
			maxPathLength = len(path)
			maxPathExample = path
		}

		if info.IsDir() {
			return nil
		}

		size += info.Size()

		return nil
	})
	if err != nil {
		return nil, err
	}

	info, err := e.Info()
	if err != nil {
		return nil, err
	}

	return &entryMeta{
		Name:           info.Name(),
		Mode:           info.Mode(),
		MaxPathLength:  maxPathLength,
		MaxPathExample: maxPathExample,
		Size:           size,
		NFiles:         nFiles,
	}, nil
}

func lsCmd() error {
	root := "./"

	dirEntries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	var output []*entryMeta
	for _, entry := range dirEntries {
		outputEntry, err := getEntryMeta(root, entry)
		if err != nil {
			return err
		}
		output = append(output, outputEntry)
	}

	if sortByNFiles {
		sort.Slice(output, func(i, j int) bool {
			return output[i].NFiles > output[j].NFiles
		})
	}

	if sortByPathLength {
		sort.Slice(output, func(i, j int) bool {
			return output[i].MaxPathLength > output[j].MaxPathLength
		})
	}

	if sortBySize {
		sort.Slice(output, func(i, j int) bool {
			return output[i].Size > output[j].Size
		})
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Permissions", "Name", "Max path length", "Size", "# files")
	if displayLongestPath {
		tbl = table.New("Permissions", "Name", "Max path length", "Size", "# files", "Max path example")
	}
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, outputEntry := range output {
		tbl.AddRow(outputEntry.Mode, outputEntry.Name, outputEntry.MaxPathLength, misc.PrettifyByteSize(outputEntry.Size), outputEntry.NFiles, outputEntry.MaxPathExample)
	}

	tbl.Print()
	return nil
}
