package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var directory string

var DirHashCmd = &cobra.Command{
	Use:   "dirhash",
	Short: "utils for recursive directory hashing and verification",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dirhashCmd()
	},
}

func init() {
	DirHashCmd.Flags().StringVarP(
		&directory,
		"directory",
		"d",
		"",
		"directory to recursively hash",
	)
}

func getFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func dirhashCmd() error {
	if directory == "" {
		return errors.New("directory is required")
	}

	hashByPath := make(map[string]string)
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}

		hash, err := getFileHash(path)
		if err != nil {
			return err
		}

		hashByPath[relPath] = hash
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	paths := make([]string, 0, len(hashByPath))
	for path := range hashByPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	orderedHashes := make([]map[string]string, 0, len(paths))
	for _, path := range paths {
		orderedHashes = append(orderedHashes, map[string]string{
			path: hashByPath[path],
		})
	}

	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	return encoder.Encode(orderedHashes)
}
