package lifetime_verifier

import (
	"fmt"
	"github.com/galqiwi/garc/internal/verify/verifiers/verifier"
	"os"
	"path/filepath"
	"time"
)

type LifetimeVerifierConfig struct {
	Dirs           []string `yaml:"dirs"`
	WarnLifetimeS  int64    `yaml:"warn_lifetime_s"`
	ErrorLifetimeS int64    `yaml:"max_lifetime_s"`
}

func NewLifetimeVerifier(cfg *LifetimeVerifierConfig) verifier.Verifier {
	return &lifetimeVerifier{
		cfg: cfg,
	}
}

type lifetimeVerifier struct {
	cfg *LifetimeVerifierConfig
}

func (v *lifetimeVerifier) Name() string {
	return "lifetime_verifier"
}

func (v *lifetimeVerifier) Verify() []verifier.ErrorMessage {
	var output []verifier.ErrorMessage

	for _, dirPath := range v.cfg.Dirs {
		messages, err := v.verifyDir(dirPath)
		if err != nil {
			return []verifier.ErrorMessage{verifier.NewErrorMessageFromError(err)}
		}
		output = append(output, messages...)
	}

	return output
}

func (v *lifetimeVerifier) verifyDir(dirPath string) ([]verifier.ErrorMessage, error) {
	dirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, err
	}

	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var output []verifier.ErrorMessage

	for _, entry := range dirEntries {
		entry, err := entry.Info()
		if err != nil {
			return nil, err
		}
		fileLifetime := time.Now().Sub(entry.ModTime()).Round(time.Second)
		errLifetimeLimit := time.Second * time.Duration(v.cfg.ErrorLifetimeS)
		warnLifetimeLimit := time.Second * time.Duration(v.cfg.WarnLifetimeS)
		if fileLifetime < warnLifetimeLimit && fileLifetime < errLifetimeLimit {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		fileLifetimeDays := float64(fileLifetime) / float64(time.Hour*24)
		errorLevel := verifier.Warn
		if fileLifetime > errLifetimeLimit {
			errorLevel = verifier.Error
		}

		output = append(output, verifier.ErrorMessage{
			Path:   fullPath,
			Reason: fmt.Sprintf("too old (%.0f days)", fileLifetimeDays),
			Level:  errorLevel,
		})
	}

	return output, nil
}
