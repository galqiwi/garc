package ssh_utils

import (
	"fmt"
	"github.com/galqiwi/garc/internal/utils/ssh"
	"io"
	"strconv"
	"strings"
)

func CreateRemoteFile(host ssh.Host, username string, path string, contentSource io.Reader) error {
	err := validatePath(path)
	if err != nil {
		return nil
	}
	cmd := host.Command(username, fmt.Sprintf("cp /dev/stdin %v", path))
	cmd.SetStdin(contentSource)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ReadRemoteFile(host ssh.Host, username string, path string, to io.Writer) error {
	err := validatePath(path)
	if err != nil {
		return nil
	}
	cmd := host.Command(username, fmt.Sprintf("cp %v /dev/stdout", path))
	cmd.SetStdout(to)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func GetFileSize(host ssh.Host, username string, path string) (int64, error) {
	err := validatePath(path)
	if err != nil {
		return 0, err
	}
	outputStr, err := RunAndGetStdout(host.Command(
		username,
		fmt.Sprintf("stat -c %%s %v", path),
	))
	if err != nil {
		return 0, err
	}

	outputStr = strings.Trim(outputStr, "\n")

	output, err := strconv.ParseInt(outputStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return output, nil
}
