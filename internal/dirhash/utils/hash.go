package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/sync/errgroup"
)

type HashMeta map[string]string

func GetHashMeta(dirpath string) (HashMeta, error) {
	paths, err := getFilePaths(dirpath)
	if err != nil {
		return nil, err
	}

	output := make(HashMeta)
	var mu sync.Mutex

	g := new(errgroup.Group)

	for _, path := range paths {
		path := path // capture loop variable
		g.Go(func() error {
			relPath, err := filepath.Rel(dirpath, path)
			if err != nil {
				return err
			}

			hash, err := getFileHash(path)
			if err != nil {
				return err
			}

			mu.Lock()
			output[relPath] = hash
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return output, nil
}

// Returns all files in the directory and all its subdirectories
func getFilePaths(dirpath string) ([]string, error) {
	output := make([]string, 0)

	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		output = append(output, path)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return output, nil
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
