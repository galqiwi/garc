package ssh_utils

import (
	"fmt"
	"github.com/galqiwi/garc/internal/utils/ssh"
	"strings"
)

func DoesDirExist(host ssh.Host, username string, dirPath string) (bool, error) {
	err := validatePath(dirPath)
	if err != nil {
		return false, err
	}
	output, err := RunAndGetStdout(host.Command(
		username,
		fmt.Sprintf("[ -d %v ] && echo 1 || echo 0", dirPath),
	))
	if err != nil {
		return false, err
	}
	if output == "1\n" {
		return true, nil
	}
	if output == "0\n" {
		return false, nil
	}
	return false, fmt.Errorf("invalid remote command output: %v", output)
}

func CreateRemoteDir(host ssh.Host, username string, dirPath string) error {
	err := validatePath(dirPath)
	if err != nil {
		return err
	}
	_, err = RunAndGetStdout(host.Command(
		username,
		fmt.Sprintf("mkdir -p %v", dirPath),
	))
	if err != nil {
		return err
	}
	return nil
}

func RemoveRemoteDir(host ssh.Host, username string, dirPath string) error {
	err := validatePath(dirPath)
	if err != nil {
		return err
	}
	_, err = RunAndGetStdout(host.Command(
		username,
		fmt.Sprintf("rm -rf %v", dirPath),
	))
	if err != nil {
		return err
	}
	return nil
}

func ListRemoteDir(host ssh.Host, username string, dirPath string) ([]string, error) {
	err := validatePath(dirPath)
	if err != nil {
		return nil, err
	}
	lsOutput, err := RunAndGetStdout(host.Command(
		username,
		fmt.Sprintf("ls -1a %v", dirPath),
	))
	if err != nil {
		return nil, err
	}

	var output []string
	for _, path := range strings.Split(lsOutput, "\n") {
		if path == "." || path == ".." || path == "" {
			continue
		}
		output = append(output, path)
	}

	return output, nil
}
