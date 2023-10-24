package common

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

type ArchiveMeta struct {
	Id               uuid.UUID
	Name             string
	CreationTime     time.Time
	ModificationTime time.Time
	Version          int64
}

func tryGetHostname() string {
	output, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return output
}

func NewArchiveMeta(name string) *ArchiveMeta {
	hostname := tryGetHostname()
	creationTime := time.Now()
	if name == "" {
		name = fmt.Sprintf("%v (%v)", hostname, creationTime.Format("2006-01-02.15:04:05"))
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
