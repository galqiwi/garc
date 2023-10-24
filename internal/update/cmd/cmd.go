package cmd

import (
	"errors"
	"fmt"
	"github.com/galqiwi/garc/internal/config"
	"github.com/galqiwi/garc/internal/utils/misc"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

var UnsetSourceErr = errors.New("source is not set, run garc config edit")

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update this binary",
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateCmd()
	},
}

func updateCmd() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}
	if cfg.UpdateConfig.SourceURL == "" {
		return UnsetSourceErr
	}
	currentExec, err := misc.GetCurrentExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %e", err)
	}
	_ = currentExec
	res, err := http.Get(cfg.UpdateConfig.SourceURL)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	file, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(file.Name())
	}()
	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	err = file.Chmod(0700)
	if err != nil {
		return err
	}

	err = os.Rename(file.Name(), currentExec)
	if err != nil {
		return err
	}

	fmt.Println("successfully updated itself")

	return nil
}
