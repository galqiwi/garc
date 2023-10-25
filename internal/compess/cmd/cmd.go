package cmd

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/galqiwi/garc/internal/utils/tarball"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"path/filepath"
)

var CompressCmd = &cobra.Command{
	Use:   "compress [flags] dirname",
	Short: "tar czf with modification time preservation",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Help()
		}
		return compressCmd(args[0])
	},
}

func init() {
	CompressCmd.AddCommand()
}

func isDir(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = file.Close()
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func getTarballPath(path string) string {
	return filepath.Join(
		filepath.Dir(path),
		filepath.Base(path)+".tar.gz",
	)
}

func validatePath(path string) error {
	pathIsDir, err := isDir(path)
	if err != nil {
		return err
	}
	if !pathIsDir {
		return fmt.Errorf("%v is not a directory", path)
	}

	if path == filepath.Dir(path) {
		return fmt.Errorf("path %v should have parent", path)
	}

	tarballPath := getTarballPath(path)

	if _, err = os.Stat(tarballPath); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%v already exists", tarballPath)
	}

	return nil
}

func compressCmd(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	err = validatePath(path)
	if err != nil {
		return err
	}

	eg := errgroup.Group{}

	tarReader, tarWriter := io.Pipe()

	gzipTarWriter := gzip.NewWriter(tarWriter)

	eg.Go(func() error {
		defer func() {
			_ = gzipTarWriter.Close()
			_ = tarWriter.Close()
		}()
		return tarball.CreateTarball(path, filepath.Base(path), gzipTarWriter, nil)
	})

	eg.Go(func() error {
		file, err := os.Create(getTarballPath(path))
		if err != nil {
			return err
		}
		defer func() {
			_ = file.Close()
		}()

		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}

		return nil
	})

	err = eg.Wait()
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if err := os.Chtimes(getTarballPath(path), fileInfo.ModTime(), fileInfo.ModTime()); err != nil {
		return fmt.Errorf("error setting timestamps for file %s: %v", getTarballPath(path), err)
	}

	return nil
}
