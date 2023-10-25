package tarball

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func setByList(list []string) map[string]struct{} {
	output := make(map[string]struct{})
	for _, elem := range list {
		output[elem] = struct{}{}
	}
	return output
}

func CreateTarball(rootPath string, innerPath string, buf io.Writer, toExclude []string) error {
	tw := tar.NewWriter(buf)
	defer tw.Close()

	excludedPaths := setByList(toExclude)

	return filepath.Walk(rootPath, func(fullFilePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking to %s: %w", fullFilePath, err)
		}

		if _, isExcluded := excludedPaths[fullFilePath]; isExcluded {
			return nil
		}

		relFilePath, err := filepath.Rel(rootPath, fullFilePath)
		if err != nil {
			return fmt.Errorf("error getting relative path for %s: %w", fullFilePath, err)
		}

		filePath := relFilePath
		if innerPath != "" {
			filePath = filepath.Join(innerPath, filePath)
		}
		header, err := tar.FileInfoHeader(fileInfo, filePath)
		if err != nil {
			return fmt.Errorf("error creating tar header for %s: %w", relFilePath, err)
		}

		header.Name = filepath.ToSlash(filePath)
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("error writing tar header for %s: %w", relFilePath, err)
		}

		if !fileInfo.IsDir() {
			data, err := os.Open(fullFilePath)
			if err != nil {
				return fmt.Errorf("error opening file %s: %w", fullFilePath, err)
			}
			defer func() {
				_ = data.Close()
			}()

			if _, err := io.Copy(tw, data); err != nil {
				return fmt.Errorf("error copying data for %s: %w", fullFilePath, err)
			}
		}
		return nil
	})
}

func ExtractTarball(tarball io.Reader, destPath string) error {
	tr := tar.NewReader(tarball)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(destPath, filepath.FromSlash(header.Name))
		info := header.FileInfo()

		if info.IsDir() {
			if err := os.MkdirAll(targetPath, info.Mode()); err != nil {
				return fmt.Errorf("error creating directory %s: %v", targetPath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("error creating parent directory for %s: %v", targetPath, err)
		}

		file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return fmt.Errorf("error opening file %s for writing: %v", targetPath, err)
		}

		if _, err := io.Copy(file, tr); err != nil {
			_ = file.Close()
			return fmt.Errorf("error writing to file %s: %v", targetPath, err)
		}
		_ = file.Close()

		if err := os.Chtimes(targetPath, header.AccessTime, header.ModTime); err != nil {
			return fmt.Errorf("error setting timestamps for file %s: %v", targetPath, err)
		}
	}

	return nil
}
