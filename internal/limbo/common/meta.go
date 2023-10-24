package common

import (
	"encoding/json"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"time"
)

type ArchiveMeta struct {
	Id               uuid.UUID
	Name             string
	CreationTime     time.Time
	ModificationTime time.Time
	Version          int64
}

func getDefaultName() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Base(wd)
}

func NewArchiveMeta(name string) *ArchiveMeta {
	creationTime := time.Now()
	if name == "" {
		name = getDefaultName()
	}
	return &ArchiveMeta{
		Id:               uuid.New(),
		Name:             name,
		CreationTime:     creationTime,
		ModificationTime: creationTime,
	}
}

func ReadCurrentArchiveMeta() (*ArchiveMeta, error) {
	file, err := os.Open(LimboMetadataFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var output ArchiveMeta
	err = json.NewDecoder(file).Decode(&output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func WriteCurrentArchiveMeta(meta *ArchiveMeta) error {
	file, err := os.Create(LimboMetadataFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	err = json.NewEncoder(file).Encode(meta)
	if err != nil {
		return err
	}

	return nil
}
