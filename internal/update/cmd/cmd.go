package cmd

import (
	"fmt"
	"github.com/galqiwi/garc/internal/utils/misc"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update this binary",
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateCmd()
	},
}

func downloadExecutable(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return file.Chmod(0700)
}

func updateCmd() error {
	url, err := getLatestReleaseURL("galqiwi", "garc", "garc")

	currentExec, err := misc.GetCurrentExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %e", err)
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	tmpPath := filepath.Join(tmpDir, "executable")

	err = downloadExecutable(url, tmpPath)
	if err != nil {
		return err
	}

	versionProbeResp, err := exec.Command(tmpPath, "version").CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"downloaded binary is invalid. Version probe responded with \n%v\n%v",
			string(versionProbeResp),
			err,
		)
	}

	err = os.Rename(tmpPath, currentExec)
	if err != nil {
		return err
	}

	fmt.Println("successfully updated itself")

	return nil
}
