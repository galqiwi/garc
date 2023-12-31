package config

import (
	"fmt"
	"github.com/galqiwi/garc/internal/limbo/common"
	"github.com/galqiwi/garc/internal/update/config"
	config2 "github.com/galqiwi/garc/internal/verify/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	LimboConfig     common.LimboConfig      `yaml:"limbo"`
	UpdateConfig    config.UpdateConfig     `yaml:"update"`
	VerifiersConfig config2.VerifiersConfig `yaml:"verifiers"`
}

var configPath string

func getDefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}
	return filepath.Join(home, ".config", "garc", "config.yaml"), nil
}

func GetConfigPath() (string, error) {
	if configPath != "" {
		return configPath, nil
	}
	return getDefaultConfigPath()
}

func AddConfigFlag(cmd *cobra.Command) {
	defaultConfigPath, err := getDefaultConfigPath()
	if err != nil {
		defaultConfigPath = "UNKNOWN"
	}
	cmd.PersistentFlags().StringVar(
		&configPath,
		"config",
		"",
		fmt.Sprintf("config file (default is %v)", defaultConfigPath),
	)
}

func GetConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var output Config
	err = yaml.NewDecoder(file).Decode(&output)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func SaveConfig(config *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}

	output, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, output, 0600)
}
